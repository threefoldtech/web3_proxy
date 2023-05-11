package goethclient

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	tft "github.com/threefoldfoundation/tft/bridge/stellar/contracts/tokenv1"
)

const (
	PublicEthTftContractAddress = "0x395E925834996e558bdeC77CD648435d620AfB5b"
)

func (c *Client) TransferTftEth(ctx context.Context, destination string, amount int64) (string, error) {
	return c.TransferTokens(ctx, PublicEthTftContractAddress, destination, amount)
}

func (c *Client) WithdrawEthTftToStellar(ctx context.Context, destination string, amount int64) (string, error) {
	tft, err := tft.NewToken(common.HexToAddress(PublicEthTftContractAddress), c.Eth)
	if err != nil {
		return "", err
	}

	nonce, err := c.Eth.PendingNonceAt(context.Background(), c.Address)
	if err != nil {
		return "", errors.Wrap(err, "failed to get nonce")
	}

	gasPrice, err := c.Eth.SuggestGasPrice(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "failed to suggest gas price")
	}

	ctxWithCancel, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	opts := &bind.TransactOpts{
		GasPrice: gasPrice,
		GasLimit: GasLimit,
		Nonce:    big.NewInt(int64(nonce)),
		Context:  ctxWithCancel,
	}

	tx, err := tft.Withdraw(opts, big.NewInt(amount), destination, "stellar")
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (c *Client) GetTftBalance(ctx context.Context) (*big.Int, error) {
	tft, err := tft.NewToken(common.HexToAddress(PublicEthTftContractAddress), c.Eth)
	if err != nil {
		return nil, err
	}

	ctxWithCancel, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return tft.BalanceOf(&bind.CallOpts{
		Context: ctxWithCancel,
	}, c.Address)
}