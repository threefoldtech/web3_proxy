module eth

import freeflowuniverse.crystallib.rpcwebsocket { RpcWsClient }

const (
	default_timeout = 500000
)

[params]
pub struct Load {
	url string
	secret string
}

[params]
pub struct Transfer {
	destination string
	amount i64
}

[params]
pub struct TokenTransfer {
	contract_address string
	destination string
	amount i64
}

[params]
pub struct TokenTransferFrom {
	contract_address string
	from string
	destination string
	amount i64
}

[params]
pub struct ApproveTokenSpending {
	contract_address string
	target          string
	amount          i64
}

[params]
pub struct MultisigOwner {
	contract_address string
	target string
	threshold i64
}

[params]
pub struct ApproveHash {
	contract_address string
	hash string
}

[params]
pub struct InitiateMultisigEthTransfer  {
	contract_address string
	destination string
	amount i64
}

[params]
pub struct InitiateMultisigTokenTransfer {
	contract_address string
	token_address string
	destination string
	amount i64
}

[params]
pub struct GetFungibleBalance{
	contract_address string
	target          string
}

[params]
pub struct OwnerOfFungible {
	contract_address string
	token_id i64
}

[params]
pub struct TransferFungible{
	contract_address string
	from string
	to string
	token_id i64
}

[params]
pub struct SetFungibleApproval{
	contract_address string
	from string
	to string
	amount i64
}

[params]
pub struct SetFungibleApprovalForAll{
	contract_address string
	from string
	to string
	approved bool
}

[params]
pub struct ApprovalForFungible{
	contract_address string
	owner string
	operator string
}

[noinit]
pub struct EthClient {
mut:
	client &RpcWsClient
}

pub fn new(mut client RpcWsClient) EthClient {
	return EthClient{
		client: &client
	}
}

// CORE

pub fn (mut e EthClient) load(args Load) ! {
	_ := e.client.send_json_rpc[[]Load, string]('eth.Load', [args], eth.default_timeout)!
}

pub fn (mut e EthClient) transer(args Transfer) !string {
	return e.client.send_json_rpc[[]Transfer, string]('eth.Transfer', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) balance(address string) !i64 {
	return e.client.send_json_rpc[[]string, i64]('eth.Balance', [address], eth.default_timeout)!
}

pub fn (mut e EthClient) height() !u64 {
	return e.client.send_json_rpc[[]string, u64]('eth.Height', []string{}, eth.default_timeout)!
}

// ERC20

pub fn (mut e EthClient) token_balance(contractAddress string, address string) !i64 {
	return e.client.send_json_rpc[[]string, i64]('eth.GetTokenBalance', [contractAddress, address], eth.default_timeout)!
}

pub fn (mut e EthClient) token_transer(args TokenTransfer) !string {
	return e.client.send_json_rpc[[]TokenTransfer, string]('eth.TransferTokens', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) token_transer_from(args TokenTransferFrom) !string {
	return e.client.send_json_rpc[[]TokenTransferFrom, string]('eth.TransferFromTokens', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) approve_token_spending(args ApproveTokenSpending) !string {
	return e.client.send_json_rpc[[]ApproveTokenSpending, string]('eth.ApproveTokenSpending', [args],
		eth.default_timeout)!
}

// Multisig

pub fn (mut e EthClient) get_multisig_owners(contractAddress string) ![]string {
	return e.client.send_json_rpc[[]string, []string]('eth.GetMultisigOwners', [contractAddress],
		eth.default_timeout)!
}

pub fn (mut e EthClient) get_multisig_threshold(contractAddress string) !i64 {
	return e.client.send_json_rpc[[]string, i64]('eth.GetMultisigThreshold', [contractAddress],
		eth.default_timeout)!
}

pub fn (mut e EthClient) add_multisig_owner(args MultisigOwner) !string {
	return e.client.send_json_rpc[[]MultisigOwner, string]('eth.AddMultisigOwner', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) remove_multisig_owner(args MultisigOwner) !string {
	return e.client.send_json_rpc[[]MultisigOwner, string]('eth.RemoveMultisigOwner', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) approve_hash(args ApproveHash) !string {
	return e.client.send_json_rpc[[]ApproveHash, string]('eth.ApproveHash', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) is_approved(args ApproveHash) !bool {
	return e.client.send_json_rpc[[]ApproveHash, bool]('eth.IsApproved', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) initiate_multisig_eth_transfer(args InitiateMultisigEthTransfer) !string {
	return e.client.send_json_rpc[[]InitiateMultisigEthTransfer, string]('eth.InitiateMultisigEthTransfer', [args],
		eth.default_timeout)!
}

pub fn (mut e EthClient) initiate_multisig_token_transfer(args InitiateMultisigTokenTransfer) !string {
	return e.client.send_json_rpc[[]InitiateMultisigTokenTransfer, string]('eth.InitiateMultisigTokenTransfer', [args],
		eth.default_timeout)!
}

// Fungibles

// GetFungibleBalance (ERC721)
pub fn (mut e EthClient) get_fungible_balance(args GetFungibleBalance) !i64 {
	return e.client.send_json_rpc[[]GetFungibleBalance, i64]('eth.GetFungibleBalance', [args],
		eth.default_timeout)!
}

// OwnerOfFungible (ERC721)
pub fn (mut e EthClient) owner_of_fungible(args OwnerOfFungible) !string {
	return e.client.send_json_rpc[[]OwnerOfFungible, string]('eth.OwnerOfFungible', [args],
		eth.default_timeout)!
}

// SafeTransferFungible (ERC721)
pub fn (mut e EthClient) safe_transfer_fungible(args TransferFungible) !string {
	return e.client.send_json_rpc[[]TransferFungible, string]('eth.SafeTransferFungible', [args],
		eth.default_timeout)!
}

// TransferFungible (ERC721)
pub fn (mut e EthClient) transfer_fungible(args TransferFungible) !string {
	return e.client.send_json_rpc[[]TransferFungible, string]('eth.TransferFungible', [args],
		eth.default_timeout)!
}

// SetFungibleApproval (ERC721)
pub fn (mut e EthClient) set_fungible_approval(args SetFungibleApproval) !string {
	return e.client.send_json_rpc[[]SetFungibleApproval, string]('eth.SetFungibleApproval', [args],
		eth.default_timeout)!
}

// SetFungibleApprovalForAll (ERC721)
pub fn (mut e EthClient) set_fungible_approval_for_all(args SetFungibleApprovalForAll) !string {
	return e.client.send_json_rpc[[]SetFungibleApprovalForAll, string]('eth.SetFungibleApprovalForAll', [args],
		eth.default_timeout)!
}

// GetApprovalForFungible (ERC721)
pub fn (mut e EthClient) get_fungible_approval(args ApprovalForFungible) !bool {
	return e.client.send_json_rpc[[]ApprovalForFungible, bool]('eth.GetApprovalForFungible', [args],
		eth.default_timeout)!
}

// GetApprovalForAllFungible (ERC721)
pub fn (mut e EthClient) get_fungible_approval_for_all(args ApprovalForFungible) !bool {
	return e.client.send_json_rpc[[]ApprovalForFungible, bool]('eth.GetApprovalForAllFungible', [args],
		eth.default_timeout)!
}