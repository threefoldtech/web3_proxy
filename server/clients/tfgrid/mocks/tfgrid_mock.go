// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/server/procedures/tfgrid.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	graphql "github.com/threefoldtech/grid3-go/graphql"
	client "github.com/threefoldtech/grid3-go/node"
	workloads "github.com/threefoldtech/grid3-go/workloads"
)

// MockTFGridClient is a mock of TFGridClient interface.
type MockTFGridClient struct {
	ctrl     *gomock.Controller
	recorder *MockTFGridClientMockRecorder
}

// MockTFGridClientMockRecorder is the mock recorder for MockTFGridClient.
type MockTFGridClientMockRecorder struct {
	mock *MockTFGridClient
}

// NewMockTFGridClient creates a new mock instance.
func NewMockTFGridClient(ctrl *gomock.Controller) *MockTFGridClient {
	mock := &MockTFGridClient{ctrl: ctrl}
	mock.recorder = &MockTFGridClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTFGridClient) EXPECT() *MockTFGridClientMockRecorder {
	return m.recorder
}

// CancelProject mocks base method.
func (m *MockTFGridClient) CancelProject(ctx context.Context, projectName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelProject", ctx, projectName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelProject indicates an expected call of CancelProject.
func (mr *MockTFGridClientMockRecorder) CancelProject(ctx, projectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelProject", reflect.TypeOf((*MockTFGridClient)(nil).CancelProject), ctx, projectName)
}

// DeployDeployment mocks base method.
func (m *MockTFGridClient) DeployDeployment(ctx context.Context, d *workloads.Deployment) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployDeployment", ctx, d)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeployDeployment indicates an expected call of DeployDeployment.
func (mr *MockTFGridClientMockRecorder) DeployDeployment(ctx, d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployDeployment", reflect.TypeOf((*MockTFGridClient)(nil).DeployDeployment), ctx, d)
}

// DeployGWFQDN mocks base method.
func (m *MockTFGridClient) DeployGWFQDN(ctx context.Context, gw *workloads.GatewayFQDNProxy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployGWFQDN", ctx, gw)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeployGWFQDN indicates an expected call of DeployGWFQDN.
func (mr *MockTFGridClientMockRecorder) DeployGWFQDN(ctx, gw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployGWFQDN", reflect.TypeOf((*MockTFGridClient)(nil).DeployGWFQDN), ctx, gw)
}

// DeployGWName mocks base method.
func (m *MockTFGridClient) DeployGWName(ctx context.Context, gw *workloads.GatewayNameProxy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployGWName", ctx, gw)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeployGWName indicates an expected call of DeployGWName.
func (mr *MockTFGridClientMockRecorder) DeployGWName(ctx, gw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployGWName", reflect.TypeOf((*MockTFGridClient)(nil).DeployGWName), ctx, gw)
}

// DeployK8sCluster mocks base method.
func (m *MockTFGridClient) DeployK8sCluster(ctx context.Context, k8s *workloads.K8sCluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployK8sCluster", ctx, k8s)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeployK8sCluster indicates an expected call of DeployK8sCluster.
func (mr *MockTFGridClientMockRecorder) DeployK8sCluster(ctx, k8s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployK8sCluster", reflect.TypeOf((*MockTFGridClient)(nil).DeployK8sCluster), ctx, k8s)
}

// DeployNetwork mocks base method.
func (m *MockTFGridClient) DeployNetwork(ctx context.Context, znet *workloads.ZNet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployNetwork", ctx, znet)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeployNetwork indicates an expected call of DeployNetwork.
func (mr *MockTFGridClientMockRecorder) DeployNetwork(ctx, znet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployNetwork", reflect.TypeOf((*MockTFGridClient)(nil).DeployNetwork), ctx, znet)
}

// GetNodeClient mocks base method.
func (m *MockTFGridClient) GetNodeClient(nodeID uint32) (*client.NodeClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeClient", nodeID)
	ret0, _ := ret[0].(*client.NodeClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeClient indicates an expected call of GetNodeClient.
func (mr *MockTFGridClientMockRecorder) GetNodeClient(nodeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeClient", reflect.TypeOf((*MockTFGridClient)(nil).GetNodeClient), nodeID)
}

// GetProjectContracts mocks base method.
func (m *MockTFGridClient) GetProjectContracts(ctx context.Context, projectName string) (graphql.Contracts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectContracts", ctx, projectName)
	ret0, _ := ret[0].(graphql.Contracts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectContracts indicates an expected call of GetProjectContracts.
func (mr *MockTFGridClientMockRecorder) GetProjectContracts(ctx, projectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectContracts", reflect.TypeOf((*MockTFGridClient)(nil).GetProjectContracts), ctx, projectName)
}
