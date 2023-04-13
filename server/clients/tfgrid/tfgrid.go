package tfgrid

import (
	"context"

	"github.com/pkg/errors"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/graphql"
	client "github.com/threefoldtech/tfgrid-sdk-go/grid-client/node"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

type TFGridClient interface {
	DeployGWFQDN(ctx context.Context, gw *workloads.GatewayFQDNProxy) error
	DeployGWName(ctx context.Context, gw *workloads.GatewayNameProxy) error
	DeployK8sCluster(ctx context.Context, k8s *workloads.K8sCluster) error
	DeployNetwork(ctx context.Context, znet *workloads.ZNet) error
	DeployDeployment(ctx context.Context, d *workloads.Deployment) (uint64, error)
	CancelProject(ctx context.Context, projectName string) error
	GetProjectContracts(ctx context.Context, projectName string) (graphql.Contracts, error)
	GetNodeClient(nodeID uint32) (*client.NodeClient, error)
	RMBCall(ctx context.Context, twin uint32, fn string, data interface{}, result interface{}) error
	FilterNodes(filter types.NodeFilter, pagination types.Limit) (res []types.Node, totalCount int, err error)
	FilterFarms(filter types.FarmFilter, pagination types.Limit) (res []types.Farm, totalCount int, err error)
	GetNode(nodeID uint32) (res types.NodeWithNestedCapacity, err error)
}

type tfgridClient struct {
	client *deployer.TFPluginClient
}

func NewTFGridClient(c *deployer.TFPluginClient) TFGridClient {
	return &tfgridClient{
		client: c,
	}
}

func (c *tfgridClient) DeployGWFQDN(ctx context.Context, gw *workloads.GatewayFQDNProxy) error {
	if err := c.client.GatewayFQDNDeployer.Deploy(ctx, gw); err != nil {
		return errors.Wrapf(err, "failed to deploy gateway fqdn")
	}

	return nil
}

func (c *tfgridClient) DeployGWName(ctx context.Context, gw *workloads.GatewayNameProxy) error {
	if err := c.client.GatewayNameDeployer.Deploy(ctx, gw); err != nil {
		return errors.Wrapf(err, "failed to deploy gateway %s", gw.Name)
	}

	return nil
}

func (c *tfgridClient) DeployK8sCluster(ctx context.Context, k8s *workloads.K8sCluster) error {
	if err := c.client.K8sDeployer.Deploy(ctx, k8s); err != nil {
		return errors.Wrapf(err, "Failed to deploy K8s Cluster")
	}

	return nil
}

func (c *tfgridClient) DeployNetwork(ctx context.Context, znet *workloads.ZNet) error {
	if err := c.client.NetworkDeployer.Deploy(ctx, znet); err != nil {
		return errors.Wrap(err, "failed to deploy network")
	}

	return nil
}
func (c *tfgridClient) DeployDeployment(ctx context.Context, d *workloads.Deployment) (uint64, error) {
	if err := c.client.DeploymentDeployer.Deploy(ctx, d); err != nil {
		return 0, errors.Wrap(err, "failed to deploy deployment")
	}

	return d.ContractID, nil
}
func (c *tfgridClient) CancelProject(ctx context.Context, projectName string) error {
	if err := c.client.CancelByProjectName(projectName); err != nil {
		return errors.Wrapf(err, "failed to cancel project %s", projectName)
	}

	return nil
}

func (c *tfgridClient) GetNodeClient(nodeID uint32) (*client.NodeClient, error) {
	nodeClient, err := c.client.NcPool.GetNodeClient(c.client.SubstrateConn, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %d client", nodeID)
	}

	return nodeClient, nil
}
func (c *tfgridClient) GetProjectContracts(ctx context.Context, projectName string) (graphql.Contracts, error) {
	contracts, err := c.client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return graphql.Contracts{}, errors.Wrapf(err, "failed to get project (%s) contracts", projectName)
	}

	return contracts, nil
}

func (c *tfgridClient) RMBCall(ctx context.Context, twin uint32, fn string, data interface{}, result interface{}) error {
	return c.client.RMB.Call(ctx, twin, fn, data, result)
}
func (c *tfgridClient) FilterNodes(filter types.NodeFilter, pagination types.Limit) (res []types.Node, totalCount int, err error) {
	return c.client.GridProxyClient.Nodes(filter, pagination)
}
func (c *tfgridClient) FilterFarms(filter types.FarmFilter, pagination types.Limit) (res []types.Farm, totalCount int, err error) {
	return c.client.GridProxyClient.Farms(filter, pagination)
}
func (c *tfgridClient) GetNode(nodeID uint32) (res types.NodeWithNestedCapacity, err error) {
	return c.client.GridProxyClient.Node(nodeID)
}
