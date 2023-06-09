package explorer

import proxyTypes "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"

type NodesRequestParams struct {
	Filters    proxyTypes.NodeFilter `json:"filters"`
	Pagination proxyTypes.Limit      `json:"pagination"`
}

type FarmsRequestParams struct {
	Filters    proxyTypes.FarmFilter `json:"filters"`
	Pagination proxyTypes.Limit      `json:"pagination"`
}

type TwinsRequestParams struct {
	Filters    proxyTypes.TwinFilter `json:"filters"`
	Pagination proxyTypes.Limit      `json:"pagination"`
}

type ContractsRequestParams struct {
	Filters    proxyTypes.ContractFilter `json:"filters"`
	Pagination proxyTypes.Limit          `json:"pagination"`
}

type NodesResult struct {
	Nodes      []proxyTypes.Node `json:"nodes"`
	TotalCount int               `json:"total_count"`
}

type FarmsResult struct {
	Farms      []proxyTypes.Farm `json:"farms"`
	TotalCount int               `json:"total_count"`
}
type TwinsResult struct {
	Twins      []proxyTypes.Twin `json:"twins"`
	TotalCount int               `json:"total_count"`
}
type ContractsResult struct {
	Contracts  []proxyTypes.Contract `json:"contracts"`
	TotalCount int                   `json:"total_count"`
}
