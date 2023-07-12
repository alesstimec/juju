// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/secrets (interfaces: JujuAPIClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	secrets "github.com/juju/juju/core/secrets"
	secrets0 "github.com/juju/juju/secrets"
	provider "github.com/juju/juju/secrets/provider"
	gomock "go.uber.org/mock/gomock"
)

// MockJujuAPIClient is a mock of JujuAPIClient interface.
type MockJujuAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockJujuAPIClientMockRecorder
}

// MockJujuAPIClientMockRecorder is the mock recorder for MockJujuAPIClient.
type MockJujuAPIClientMockRecorder struct {
	mock *MockJujuAPIClient
}

// NewMockJujuAPIClient creates a new mock instance.
func NewMockJujuAPIClient(ctrl *gomock.Controller) *MockJujuAPIClient {
	mock := &MockJujuAPIClient{ctrl: ctrl}
	mock.recorder = &MockJujuAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJujuAPIClient) EXPECT() *MockJujuAPIClientMockRecorder {
	return m.recorder
}

// GetBackendConfigForDrain mocks base method.
func (m *MockJujuAPIClient) GetBackendConfigForDrain(arg0 *string) (*provider.ModelBackendConfig, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackendConfigForDrain", arg0)
	ret0, _ := ret[0].(*provider.ModelBackendConfig)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBackendConfigForDrain indicates an expected call of GetBackendConfigForDrain.
func (mr *MockJujuAPIClientMockRecorder) GetBackendConfigForDrain(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackendConfigForDrain", reflect.TypeOf((*MockJujuAPIClient)(nil).GetBackendConfigForDrain), arg0)
}

// GetContentInfo mocks base method.
func (m *MockJujuAPIClient) GetContentInfo(arg0 *secrets.URI, arg1 string, arg2, arg3 bool) (*secrets0.ContentParams, *provider.ModelBackendConfig, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentInfo", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*secrets0.ContentParams)
	ret1, _ := ret[1].(*provider.ModelBackendConfig)
	ret2, _ := ret[2].(bool)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetContentInfo indicates an expected call of GetContentInfo.
func (mr *MockJujuAPIClientMockRecorder) GetContentInfo(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentInfo", reflect.TypeOf((*MockJujuAPIClient)(nil).GetContentInfo), arg0, arg1, arg2, arg3)
}

// GetRevisionContentInfo mocks base method.
func (m *MockJujuAPIClient) GetRevisionContentInfo(arg0 *secrets.URI, arg1 int, arg2 bool) (*secrets0.ContentParams, *provider.ModelBackendConfig, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevisionContentInfo", arg0, arg1, arg2)
	ret0, _ := ret[0].(*secrets0.ContentParams)
	ret1, _ := ret[1].(*provider.ModelBackendConfig)
	ret2, _ := ret[2].(bool)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetRevisionContentInfo indicates an expected call of GetRevisionContentInfo.
func (mr *MockJujuAPIClientMockRecorder) GetRevisionContentInfo(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevisionContentInfo", reflect.TypeOf((*MockJujuAPIClient)(nil).GetRevisionContentInfo), arg0, arg1, arg2)
}

// GetSecretBackendConfig mocks base method.
func (m *MockJujuAPIClient) GetSecretBackendConfig(arg0 *string) (*provider.ModelBackendConfigInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretBackendConfig", arg0)
	ret0, _ := ret[0].(*provider.ModelBackendConfigInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretBackendConfig indicates an expected call of GetSecretBackendConfig.
func (mr *MockJujuAPIClientMockRecorder) GetSecretBackendConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretBackendConfig", reflect.TypeOf((*MockJujuAPIClient)(nil).GetSecretBackendConfig), arg0)
}
