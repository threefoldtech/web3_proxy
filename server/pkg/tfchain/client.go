package tfchain

import (
	"context"
	"errors"
	"math/big"

	"github.com/LeeSmet/go-jsonrpc"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	substrate "github.com/threefoldtech/tfchain/clients/tfchain-client-go"
	"github.com/threefoldtech/web3_proxy/server/pkg"
	"github.com/threefoldtech/web3_proxy/server/pkg/state"
)

const (
	TfchainID = "tfchain"

	tfchainMainnet = "wss://tfchain.grid.tf"
	tfchainTestnet = "wss://tfchain.test.grid.tf"
	tfchainQanet   = "wss://tfchain.qa.grid.tf"
	tfchainDevnet  = "wss://tfchain.dev.grid.tf"
)

type (
	// ErrUnknownNetwork indicates a client was requested for an unknown network
	ErrUnknownNetwork struct{}
	// Client exposing stellar methods
	Client struct {
		state *state.StateManager[TfchainState]
	}
	TfchainState struct {
		client   *substrate.Substrate
		identity substrate.Identity
		network  string
	}

	Load struct {
		Network  string `json:"network"`
		Mnemonic string `json:"mnemonic"`
	}

	Transfer struct {
		Amount      uint64 `json:"amount"`
		Destination string `json:"destination"`
	}

	CreateTwin struct {
		Relay string `json:"relay"`
		Pk    []byte `json:"pk"`
	}

	AcceptTermsAndConditions struct {
		Link string `json:"link"`
		Hash string `json:"hash"`
	}

	GetContractWithHash struct {
		NodeID uint32            `json:"node_id"`
		Hash   substrate.HexHash `json:"hash"`
	}

	CreateNodeContract struct {
		NodeID             uint32  `json:"node_id"`
		Body               string  `json:"body"`
		Hash               string  `json:"hash"`
		PublicIPs          uint32  `json:"public_ips"`
		SolutionProviderID *uint64 `json:"solution_provider_id"`
	}

	CreateRentContract struct {
		NodeID             uint32  `json:"node_id"`
		SolutionProviderID *uint64 `json:"solution_provider_id"`
	}

	ServiceContractCreate struct {
		Service  string `json:"service"`
		Consumer string `json:"consumer"`
	}

	ServiceContractBill struct {
		ContractID     uint64 `json:"contract_id"`
		VariableAmount uint64 `json:"variable_amount"`
		Metadata       string `json:"metadata"`
	}

	SetServiceContractFees struct {
		ContractID  uint64 `json:"contract_id"`
		BaseFee     uint64 `json:"base_fee"`
		VariableFee uint64 `json:"variable_fee"`
	}

	ServiceContractSetMetadata struct {
		ContractID uint64 `json:"contract_id"`
		Metadata   string `json:"metadata"`
	}

	CreateFarm struct {
		Name      string                    `json:"name"`
		PublicIPs []substrate.PublicIPInput `json:"public_ips"`
	}

	SwapToStellar struct {
		TargetStellarAddress string   `json:"target_stellar_address"`
		Amount               *big.Int `json:"amount"`
	}
)

// Error implements the error interface
func (e ErrUnknownNetwork) Error() string {
	return "only 'public' and 'testnet' networks are supported"
}

// State from a connection. If no state is present, it is initialized
func State(conState jsonrpc.State) *TfchainState {
	raw, exists := conState[TfchainID]
	if !exists {
		ns := &TfchainState{
			client:  nil,
			network: tfchainTestnet,
		}
		conState[TfchainID] = ns
		return ns
	}
	ns, ok := raw.(*TfchainState)
	if !ok {
		// This means the invariant is violated, so panic here is ok
		panic("Invalid saved state for tfchain")
	}
	return ns
}

// NewClient creates a new Client ready for use
func NewClient() *Client {
	return &Client{
		state: state.NewStateManager[TfchainState](),
	}
}

func tfchainNetworkFromNetworkString(ntwrk string) (string, error) {
	if ntwrk == "mainnet" {
		return tfchainMainnet, nil
	} else if ntwrk == "testnet" {
		return tfchainTestnet, nil
	} else if ntwrk == "qanet" {
		return tfchainQanet, nil
	} else if ntwrk == "devnet" {
		return tfchainDevnet, nil
	}

	return "", errors.New("unsupported network")
}

// Load a client, connecting to the rpc endpoint at the given URL and loading a keypair from the given mnemonic
func (c *Client) Load(ctx context.Context, conState jsonrpc.State, args Load) error {
	url, err := tfchainNetworkFromNetworkString(args.Network)
	if err != nil {
		return err
	}

	mgr := substrate.NewManager(url)
	substrateConnection, err := mgr.Substrate()
	if err != nil {
		return err
	}

	identity, err := substrate.NewIdentityFromSr25519Phrase(args.Mnemonic)
	if err != nil {
		return err
	}
	state := State(conState)
	state.client = substrateConnection
	state.identity = identity
	state.network = args.Network

	return nil
}

func (c *Client) Address(ctx context.Context, conState jsonrpc.State) (string, error) {
	state := State(conState)
	if state.client == nil {
		return "", pkg.ErrClientNotConnected{}
	}

	return state.identity.Address(), nil
}

func (c *Client) Height(ctx context.Context, conState jsonrpc.State) (uint32, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.GetCurrentHeight()
}

// Transer an amount of TFT from the loaded account to the destination.
func (c *Client) Transfer(ctx context.Context, conState jsonrpc.State, args Transfer) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	dest, err := substrate.FromAddress(args.Destination)
	if err != nil {
		return err
	}

	return state.client.Transfer(state.identity, args.Amount, dest)
}

// Balance of an account for TFT on stellar.
func (c *Client) Balance(ctx context.Context, conState jsonrpc.State, address string) (int64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	accountId, err := substrate.FromAddress(address)
	if err != nil {
		return 0, err
	}

	balance, err := state.client.GetBalance(accountId)
	if err != nil {
		return 0, err
	}

	return balance.Free.Int64(), nil
}

func (c *Client) GetTwin(ctx context.Context, conState jsonrpc.State, id uint32) (*substrate.Twin, error) {
	state := State(conState)
	if state.client == nil {
		return nil, pkg.ErrClientNotConnected{}
	}

	return state.client.GetTwin(id)
}

func (c *Client) GetTwinByPubKey(ctx context.Context, conState jsonrpc.State, address string) (uint32, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	account, err := substrate.FromAddress(address)
	if err != nil {
		return 0, err
	}

	return state.client.GetTwinByPubKey(account.PublicKey())
}

func (c *Client) CreateTwin(ctx context.Context, conState jsonrpc.State, args CreateTwin) (uint32, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.CreateTwin(state.identity, args.Relay, args.Pk)
}

func (c *Client) AcceptTermsAndConditions(ctx context.Context, conState jsonrpc.State, args AcceptTermsAndConditions) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.AcceptTermsAndConditions(state.identity, args.Link, args.Hash)
}

func (c *Client) GetNode(ctx context.Context, conState jsonrpc.State, id uint32) (*substrate.Node, error) {
	state := State(conState)
	if state.client == nil {
		return nil, pkg.ErrClientNotConnected{}
	}

	return state.client.GetNode(id)
}

func (c *Client) GetNodes(ctx context.Context, conState jsonrpc.State, farm_id uint32) ([]uint32, error) {
	state := State(conState)
	if state.client == nil {
		return []uint32{}, pkg.ErrClientNotConnected{}
	}

	return state.client.GetNodes(farm_id)
}

func (c *Client) GetFarm(ctx context.Context, conState jsonrpc.State, id uint32) (*substrate.Farm, error) {
	state := State(conState)
	if state.client == nil {
		return nil, pkg.ErrClientNotConnected{}
	}

	return state.client.GetFarm(id)
}

func (c *Client) GetFarmByName(ctx context.Context, conState jsonrpc.State, name string) (uint32, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.GetFarmByName(name)
}

func (c *Client) CreateFarm(ctx context.Context, conState jsonrpc.State, args CreateFarm) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.CreateFarm(state.identity, args.Name, args.PublicIPs)
}

func (c *Client) GetContract(ctx context.Context, conState jsonrpc.State, contract_id uint64) (*substrate.Contract, error) {
	state := State(conState)
	if state.client == nil {
		return nil, pkg.ErrClientNotConnected{}
	}

	return state.client.GetContract(contract_id)
}

func (c *Client) GetNodeContracts(ctx context.Context, conState jsonrpc.State, node_id uint32) ([]types.U64, error) {
	state := State(conState)
	if state.client == nil {
		return []types.U64{}, pkg.ErrClientNotConnected{}
	}

	return state.client.GetNodeContracts(node_id)
}

func (c *Client) GetContractIDByNameRegistration(ctx context.Context, conState jsonrpc.State, name string) (uint64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.GetContractIDByNameRegistration(name)
}

func (c *Client) GetContractWithHash(ctx context.Context, conState jsonrpc.State, args GetContractWithHash) (uint64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.GetContractWithHash(args.NodeID, args.Hash)
}

func (c *Client) CreateNameContract(ctx context.Context, conState jsonrpc.State, name string) (uint64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.CreateNameContract(state.identity, name)
}

func (c *Client) CreateNodeContract(ctx context.Context, conState jsonrpc.State, args CreateNodeContract) (uint64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.CreateNodeContract(state.identity, args.NodeID, args.Body, args.Hash, args.PublicIPs, args.SolutionProviderID)
}

func (c *Client) CreateRentContract(ctx context.Context, conState jsonrpc.State, args CreateRentContract) (uint64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	return state.client.CreateRentContract(state.identity, args.NodeID, args.SolutionProviderID)
}

func (c *Client) ServiceContractCreate(ctx context.Context, conState jsonrpc.State, args ServiceContractCreate) (uint64, error) {
	state := State(conState)
	if state.client == nil {
		return 0, pkg.ErrClientNotConnected{}
	}

	accountIdService, err := substrate.FromAddress(args.Service)
	if err != nil {
		return 0, err
	}

	accountIdConsumer, err := substrate.FromAddress(args.Consumer)
	if err != nil {
		return 0, err
	}

	return state.client.ServiceContractCreate(state.identity, accountIdService, accountIdConsumer)
}

func (c *Client) ServiceContractApprove(ctx context.Context, conState jsonrpc.State, contract_id uint64) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.ServiceContractApprove(state.identity, contract_id)
}

func (c *Client) ServiceContractBill(ctx context.Context, conState jsonrpc.State, args ServiceContractBill) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.ServiceContractBill(state.identity, args.ContractID, args.VariableAmount, args.Metadata)
}

func (c *Client) ServiceContractCancel(ctx context.Context, conState jsonrpc.State, contract_id uint64) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.ServiceContractCancel(state.identity, contract_id)
}

func (c *Client) ServiceContractReject(ctx context.Context, conState jsonrpc.State, contract_id uint64) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.ServiceContractReject(state.identity, contract_id)
}

func (c *Client) ServiceContractSetFees(ctx context.Context, conState jsonrpc.State, args SetServiceContractFees) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.ServiceContractSetFees(state.identity, args.ContractID, args.BaseFee, args.VariableFee)
}

func (c *Client) ServiceContractSetMetadata(ctx context.Context, conState jsonrpc.State, args ServiceContractSetMetadata) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.ServiceContractSetMetadata(state.identity, args.ContractID, args.Metadata)
}

func (c *Client) CancelContract(ctx context.Context, conState jsonrpc.State, contract_id uint64) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.CancelContract(state.identity, contract_id)
}

func (c *Client) BatchCancelContract(ctx context.Context, conState jsonrpc.State, contract_ids []uint64) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.BatchCancelContract(state.identity, contract_ids)
}

func (c *Client) GetZosVersion(ctx context.Context, conState jsonrpc.State) (string, error) {
	state := State(conState)
	if state.client == nil {
		return "", pkg.ErrClientNotConnected{}
	}

	return state.client.GetZosVersion()
}

func (c *Client) SwapToStellar(ctx context.Context, conState jsonrpc.State, args SwapToStellar) error {
	state := State(conState)
	if state.client == nil {
		return pkg.ErrClientNotConnected{}
	}

	return state.client.SwapToStellar(state.identity, args.TargetStellarAddress, *args.Amount)
}
