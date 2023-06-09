package atomicswap

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/threefoldtech/atomicswap/eth"
	"github.com/threefoldtech/atomicswap/stellar"
	goethclient "github.com/threefoldtech/web3_proxy/server/clients/eth"
	"github.com/threefoldtech/web3_proxy/server/clients/nostr"
	stellargoclient "github.com/threefoldtech/web3_proxy/server/clients/stellar"
)

type (
	// Driver for atomic swaps
	Driver struct {
		nostr   *nostr.Client
		eth     *goethclient.Client
		stellar *stellargoclient.Client

		saleId string
		swapId string

		stage DriverStage

		// amount of TFT to swap, this is initialized in a sell order to the maximum available
		swapAmount uint
		// amount of other token to pay
		swapPrice uint

		msges <-chan nostr.NostrEvent

		// temporary, this nees a better way
		sct *eth.SwapContractTransactor

		secret     [32]byte
		secretHash [sha256.Size]byte
	}

	DriverStage = int

	MsgBuy struct {
		Id     string `json:"id"`
		Amount uint   `json:"amount"`
	}

	MsgAccept struct {
		Id             string         `json:"id"`
		EthAddress     common.Address `json:"ethAddress"`
		StellarAddress string         `json:"stellarAddress"`
		// Explicitly communicate accepted amount of TFT and value of 1 TFT

		// Amount of TFT to swap
		Amount uint `json:"amount"`
		// SwapPrice of 1 TFT
		SwapPrice uint `json:"swapPrice"`
	}

	MsgInitiateEth struct {
		Id                  string             `json:"id"`
		SharedSecret        [sha256.Size]byte  `json:"sharedSecret"`
		EthAddress          common.Address     `json:"ethAddress"`
		StellarAddress      string             `json:"stellarAddress"`
		InitiateTransaction *types.Transaction `json:"initiateTransaction"`
	}

	MsgParticipateStellar struct {
		Id             string `json:"id"`
		HoldingAccount string `json:"holdingAccount"`
		RefundTx       string `json:"refundTx"`
	}

	MsgRedeemed struct {
		Id     string   `json:"id"`
		Secret [32]byte `json:"secret"`
	}
)

const (
	// Initial conditions
	DriverStageOpenSale DriverStage = iota
	DriverStageStartBuy
	// In progress
	DriverStageAcceptedBuy
	DriverStageSetupSwap
	DriverStageParticipateSwap
	// Terminal conditions
	DriverStageClaimSwap
	DriverStageDone
)

const (
	// dialTimeout for dialing eth nodes
	dialTimeout = time.Second * 10
)

var (
	// contract address on the sepolia teest network
	contractAddress = common.HexToAddress("0x17f54245073bfed168a51c3d13b536e39e406063")
	// contract address on the goerli network
	// contractAddress = common.HexToAddress("0x8420c8271d602F6D0B190856Cea8E74D09A0d3cF")
	// TFT asset on the stellar testnet
	testnetTftAsset = mustStellarTestnetTftAsset()

	// goerliChainID = big.NewInt(5)
	sepoliaChainId = big.NewInt(11155111)
)

func initDriver(nostr *nostr.Client, eth *goethclient.Client, stellar *stellargoclient.Client) *Driver {
	return &Driver{
		nostr:   nostr,
		eth:     eth,
		stellar: stellar,

		swapId: uuid.NewString(),
	}
}

// Buy flow for the driver
func (d *Driver) Buy(ctx context.Context, seller string, sale nostr.Product, amount uint) error {
	d.saleId = sale.Id
	d.swapAmount = amount
	d.swapPrice = uint(sale.Price)
	d.stage = DriverStageStartBuy
	msgChan, err := d.nostr.SubscribeDirectMessagesDirect(sale.Id)
	if err != nil {
		return errors.Wrap(err, "could not subscribe to direct messages")
	}
	d.msges = msgChan

	go handleMessage(d)

	msg := MsgBuy{
		Id:     sale.Id,
		Amount: amount,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "could not encode buy message")
	}
	log.Info().Msg("Starting atomic swap buy, notify seller")
	return d.nostr.PublishDirectMessage(ctx, seller, []string{"s", sale.Id}, string(data))
}

// OpenSale on the driver
func (d *Driver) OpenSale(sale nostr.Product) error {
	d.saleId = sale.Id
	d.swapAmount = sale.Quantity
	d.swapPrice = uint(sale.Price)
	d.stage = DriverStageOpenSale
	msgChan, err := d.nostr.SubscribeDirectMessagesDirect(sale.Id)
	if err != nil {
		return errors.Wrap(err, "could not subscribe to direct messages")
	}
	d.msges = msgChan

	go handleMessage(d)

	// At this point we just wait for messages, nothing else to do here

	log.Info().Msg("Starting atomic swap sale")
	return nil
}

func (d *Driver) handleBuyMessage(ctx context.Context, sender string, req MsgBuy) {
	if req.Id != d.saleId {
		log.Debug().Msg("Ignore message which is not intended for this swap")
		return
	}
	// set swap amount
	if req.Amount > d.swapAmount {
		log.Debug().Msg("Buyer wants more TFT than we have, ignore")
		return
	}
	// Track correct buy amount for future validation
	d.swapAmount = req.Amount

	msg := MsgAccept{
		Id:             d.saleId,
		EthAddress:     d.eth.AddressFromKey(),
		StellarAddress: d.stellar.Address(),
		Amount:         d.swapAmount,
		SwapPrice:      d.swapPrice,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Can not encode accept msg")
		return
	}

	d.stage = DriverStageAcceptedBuy

	if err := d.nostr.PublishDirectMessage(ctx, sender, []string{"s", d.saleId}, string(data)); err != nil {
		log.Error().Err(err).Msg("Can not send buy accepted message")
		return
	}

	log.Info().Msg("Sent atomic swap buy offer")
}

func (d *Driver) handleBuyAcceptMessage(ctx context.Context, sender string, req MsgAccept) {
	if req.Id != d.saleId {
		log.Debug().Msg("Ignore message which is not intended for this swap")
		return
	}

	log.Info().Msg("Seller accepted our buy offer")

	// seller accepted our buy, so initiate the atomic swap.
	// we have ETH and want to buy TFT on stellar, so set up the eth
	// part of the swap
	dialCtx, cancel := context.WithTimeout(ctx, dialTimeout)
	defer cancel()
	client, err := eth.DialClient(dialCtx, d.eth.Url) // TODO: should probably be able to construct this from the existing client
	if err != nil {
		log.Error().Err(err).Msg("Failed to dial eth node")
		return
	}
	cancel()
	sct, err := eth.NewSwapContractTransactor(ctx, client, contractAddress, d.eth.Key, sepoliaChainId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to construct swap contract transactor")
		return
	}
	// save the sct so we can use it later
	d.sct = &sct
	// total wei = swap.amount * swap.price
	cost := d.swapAmount * d.swapPrice
	output, err := eth.Initiate(ctx, sct, req.EthAddress, big.NewInt(int64(cost)))
	if err != nil {
		log.Error().Err(err).Msg("Failed to initiate ETH swap")
		return
	}

	log.Info().Msgf("Submitted eth initiate transaction %s", output.ContractTransaction.Hash())

	msg := MsgInitiateEth{
		Id:                  d.saleId,
		SharedSecret:        output.SecretHash,
		EthAddress:          output.InitiatorAddress,
		StellarAddress:      d.stellar.Address(),
		InitiateTransaction: &output.ContractTransaction,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Can not encode initiate eth message")
		return
	}

	d.secretHash = output.SecretHash
	d.secret = output.Secret
	d.stage = DriverStageSetupSwap

	if err := d.nostr.PublishDirectMessage(ctx, sender, []string{"s", d.saleId}, string(data)); err != nil {
		log.Error().Err(err).Msg("Can not send buy accepted message")
		return
	}

	log.Info().Msg("Set up atomic swap on ETH side, notified seller")
}

func (d *Driver) handleInitiateEthMessage(ctx context.Context, sender string, req MsgInitiateEth) {
	if req.Id != d.saleId {
		log.Debug().Msg("Ignore message which is not intended for this swap")
		return
	}

	log.Info().Msg("Received initiate eth message")

	// Buyer initiated an Eth atomic swap, so first check and see if that is correct
	// Note that at this point, the seller does not have an sct yet
	dialCtx, cancel := context.WithTimeout(ctx, dialTimeout)
	defer cancel()
	client, err := eth.DialClient(dialCtx, d.eth.Url) // TODO: should probably be able to construct this from the existing client
	if err != nil {
		log.Error().Err(err).Msg("Failed to dial eth node")
		return
	}
	cancel()
	sct, err := eth.NewSwapContractTransactor(ctx, client, contractAddress, d.eth.Key, sepoliaChainId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to construct swap contract transactor")
		return
	}

	// save the sct so we can use it later
	d.sct = &sct

	deadline := time.Now().Add(time.Minute * 5)
	var auditOutput eth.AuditContractOutput
	for {
		auditOutput, err = eth.AuditContract(ctx, sct, req.InitiateTransaction)
		if err != nil {
			if errors.Is(err, eth.ErrTxPending) {
				if time.Now().After(deadline) {
					log.Warn().Msg("Tx not confirmed yet but deadline passed, abort")
					return
				}
				log.Info().Msg("Tx not confirmed yet, sleeping and trying again")
				time.Sleep(time.Second * 15)
				continue
			}
			log.Error().Err(err).Msg("Failed to audit eth contract")
			return
		}
		break
	}
	if err != nil {
		return
	}

	// Check the Eth locked in the contract. Notice that we will shamelessly accept if the buyer pays too much
	expectedEthValue := big.NewInt(int64(d.swapAmount * d.swapPrice))
	if auditOutput.ContractValue.Cmp(expectedEthValue) == -1 {
		log.Warn().Msg("Value in contract is less than expected value")
		return
	}

	if auditOutput.ContractAddress != contractAddress {
		log.Warn().Msg("Call is for wrong contract, ignore")
		return
	}

	if auditOutput.RecipientAddress != d.eth.AddressFromKey() {
		log.Warn().Msg("Swap is for different receiver, ignore")
		return
	}

	// TODO: shared secret doesn't have to be communicated, we can take this as the secret unilaterally
	if auditOutput.SecretHash != req.SharedSecret {
		log.Warn().Msg("Shared secret in contract call is not the same as the one communicated")
		return
	}

	// TODO: Strictly speaking we don't really care for this
	if auditOutput.RefundAddress != req.EthAddress {
		log.Warn().Msg("Contract would refund to different address than buyer")
		return
	}

	if time.Unix(auditOutput.Locktime, 0).Before(time.Now().Add(time.Hour * 2)) {
		log.Warn().Msg("Contract doesn't leave at least 2 hours to complete, ignore")
		return
	}

	// Save the secret hash for later
	d.secretHash = req.SharedSecret

	// Contract is now validated, so we participate from the stellar side
	kp := d.stellar.KeyPair()
	horizonClient := horizonclient.DefaultTestNetClient
	log.Info().Msg("Validated Eth contract, setting up stellar side")
	participateOutput, err := stellar.Participate(network.TestNetworkPassphrase, &kp, req.StellarAddress, strconv.FormatUint(uint64(d.swapAmount), 10), req.SharedSecret[:], testnetTftAsset, horizonClient)
	if err != nil {
		log.Error().Err(err).Msg("Can not participate on the stellar side")
		return
	}

	msg := MsgParticipateStellar{
		Id:             d.saleId,
		HoldingAccount: participateOutput.HoldingAccountAddress,
		RefundTx:       participateOutput.RefundTransaction,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Can not encode accept msg")
		return
	}

	d.stage = DriverStageParticipateSwap

	if err := d.nostr.PublishDirectMessage(ctx, sender, []string{"s", d.saleId}, string(data)); err != nil {
		log.Error().Err(err).Msg("Can not send buy accepted message")
		return
	}

	log.Info().Msg("Setup atomic swap on Stellar side, notified buyer")
}

func (d *Driver) handleParticipateStellarMessage(ctx context.Context, sender string, req MsgParticipateStellar) {
	if req.Id != d.saleId {
		log.Debug().Msg("Ignore message which is not intended for this swap")
		return
	}

	// Seller set up the stellar side of the swap, verify that
	horizonClient := horizonclient.DefaultTestNetClient
	refundTx := txnbuild.Transaction{}
	if err := (&refundTx).UnmarshalText([]byte(req.RefundTx)); err != nil {
		log.Warn().Err(err).Msg("Could not decode refund transaction")
		return
	}
	auditOutput, err := stellar.AuditContract(network.TestNetworkPassphrase, refundTx, req.HoldingAccount, testnetTftAsset, horizonClient)
	if err != nil {
		log.Error().Err(err).Msg("Failed to audit stellar contract")
		return
	}

	contractValue, err := strconv.ParseFloat(auditOutput.ContractValue, 64)
	if err != nil {
		log.Error().Err(err).Msg("Can not parse contract value, this is an internal coding error")
		return
	}

	// if the seller wants to give us more TFT than agreed, we will shamelessly accept
	if contractValue < float64(d.swapAmount) {
		log.Warn().Msg("Contract does not have enough TFT locked")
		return
	}

	// TODO: we loaded the "contract" based on the holding address, is this even possible?
	if auditOutput.ContractAddress != req.HoldingAccount {
		log.Warn().Msg("Contract holding address is not as specified")
		return
	}

	// Make sure we are the receiver
	if auditOutput.RecipientAddress != d.stellar.Address() {
		log.Warn().Msg("Swap is for different receiver, ignore")
		return
	}

	// Verify that the secret is properly set
	if auditOutput.SecretHash != hex.EncodeToString(d.secretHash[:]) {
		log.Warn().Msg("Shared secret in contract call is not the same as the one communicated")
		return
	}

	// TODO: Strictly speaking we don't really care for this
	//if auditOutput.RefundAddress != "TODO" {
	//	log.Warn().Msg("Contract would refund to different address than buyer")
	//	return
	//}

	if time.Unix(auditOutput.Locktime, 0).Before(time.Now().Add(time.Hour * 1)) {
		log.Warn().Msg("Contract doesn't leave at least 1 hour to complete, ignore")
		return
	}

	log.Info().Msg("Stellar contract validated, redeem it")

	// All is good in the contract, lets redeem it :)
	kp := d.stellar.KeyPair()
	redeemOutput, err := stellar.Redeem(network.TestNetworkPassphrase, &kp, req.HoldingAccount, d.secret[:], horizonClient)
	if err != nil {
		log.Error().Err(err).Msg("Failed to redeem stellar contract")
		return
	}
	log.Info().Str("Tx hash", redeemOutput.RedeemTransactionTxHash).Msg("Contract redeemed")

	// Strictly speaking it is not necessary to communicate the secret to the receiver as he can extract it from the redeem transaction
	// In fact, the seller should loop and periodically check if the contract has been redeemed, and extract the secret if it has. However
	// for now we will be courtious and notify him, while including the secret. In the end, a swap driver should continuosly build state
	// with received info to the point where he can do this on his own (this => extract secret), so that a simple "it's done" is enough
	msg := MsgRedeemed{
		Id:     d.saleId,
		Secret: d.secret,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Can not encode accept msg")
		return
	}

	d.stage = DriverStageClaimSwap

	if err := d.nostr.PublishDirectMessage(ctx, sender, []string{"s", d.saleId}, string(data)); err != nil {
		log.Error().Err(err).Msg("Can not send buy accepted message")
		return
	}

	log.Info().Msg("Redeemed atomic swap on Stellar side, notified seller, we have our TFT now")
}

func (d *Driver) handleRedeemedMessage(ctx context.Context, sender string, req MsgRedeemed) {
	if req.Id != d.saleId {
		log.Debug().Msg("Ignore message which is not intended for this swap")
		return
	}

	// So now we have a secret, for now implicitly trust the remote
	// We already validated the contract in the previous step
	redeemOutput, err := eth.Redeem(ctx, *d.sct, d.secretHash, req.Secret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to redeem eth contract")
		return
	}
	log.Info().Str("Tx hash", redeemOutput.RedeemTxHash.Hex()).Msg("Contract redeemed")

	log.Info().Msg("Redeemed atomic swap on Eth side, we have our ETH now")
}

func handleMessage(driver *Driver) {
	log.Debug().Msg("Start to handle swap messages")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for evt := range driver.msges {
		log.Warn().Int("Stage", driver.stage).Str("Sender", evt.PubKey).Msg("Got swap driver message")
		switch driver.stage {
		case DriverStageOpenSale:
			msg := MsgBuy{}
			if err := json.Unmarshal([]byte(evt.Content), &msg); err != nil {
				log.Debug().Err(err).Msg("could not decode buy message in atomic swap driver")
				continue
			}
			driver.handleBuyMessage(ctx, evt.PubKey, msg)
		case DriverStageStartBuy:
			msg := MsgAccept{}
			if err := json.Unmarshal([]byte(evt.Content), &msg); err != nil {
				log.Debug().Err(err).Msg("could not decode buy accept message in atomic swap driver")
				continue
			}
			driver.handleBuyAcceptMessage(ctx, evt.PubKey, msg)
		case DriverStageAcceptedBuy:
			msg := MsgInitiateEth{}
			if err := json.Unmarshal([]byte(evt.Content), &msg); err != nil {
				log.Debug().Err(err).Msg("could not decode initiate eth message in atomic swap driver")
				continue
			}
			driver.handleInitiateEthMessage(ctx, evt.PubKey, msg)
		case DriverStageSetupSwap:
			msg := MsgParticipateStellar{}
			if err := json.Unmarshal([]byte(evt.Content), &msg); err != nil {
				log.Debug().Err(err).Msg("could not decode participate stellar message in atomic swap driver")
				continue
			}
			driver.handleParticipateStellarMessage(ctx, evt.PubKey, msg)
		case DriverStageParticipateSwap:
			msg := MsgRedeemed{}
			if err := json.Unmarshal([]byte(evt.Content), &msg); err != nil {
				log.Debug().Err(err).Msg("could not decode redeemed message in atomic swap driver")
				continue
			}
			driver.handleRedeemedMessage(ctx, evt.PubKey, msg)
		}
	}
}

func mustStellarTestnetTftAsset() txnbuild.Asset {
	a, err := txnbuild.ParseAssetString("TFT:GA47YZA3PKFUZMPLQ3B5F2E3CJIB57TGGU7SPCQT2WAEYKN766PWIMB3")
	if err != nil {
		panic(err)
	}
	return a
}
