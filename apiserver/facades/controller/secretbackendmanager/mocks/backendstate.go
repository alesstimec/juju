// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/controller/secretbackendmanager (interfaces: BackendState)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	secrets "github.com/juju/juju/core/secrets"
	state "github.com/juju/juju/state"
	gomock "go.uber.org/mock/gomock"
)

// MockBackendState is a mock of BackendState interface.
type MockBackendState struct {
	ctrl     *gomock.Controller
	recorder *MockBackendStateMockRecorder
}

// MockBackendStateMockRecorder is the mock recorder for MockBackendState.
type MockBackendStateMockRecorder struct {
	mock *MockBackendState
}

// NewMockBackendState creates a new mock instance.
func NewMockBackendState(ctrl *gomock.Controller) *MockBackendState {
	mock := &MockBackendState{ctrl: ctrl}
	mock.recorder = &MockBackendStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackendState) EXPECT() *MockBackendStateMockRecorder {
	return m.recorder
}

// GetSecretBackendByID mocks base method.
func (m *MockBackendState) GetSecretBackendByID(arg0 string) (*secrets.SecretBackend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretBackendByID", arg0)
	ret0, _ := ret[0].(*secrets.SecretBackend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretBackendByID indicates an expected call of GetSecretBackendByID.
func (mr *MockBackendStateMockRecorder) GetSecretBackendByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretBackendByID", reflect.TypeOf((*MockBackendState)(nil).GetSecretBackendByID), arg0)
}

// SecretBackendRotated mocks base method.
func (m *MockBackendState) SecretBackendRotated(arg0 string, arg1 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretBackendRotated", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SecretBackendRotated indicates an expected call of SecretBackendRotated.
func (mr *MockBackendStateMockRecorder) SecretBackendRotated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretBackendRotated", reflect.TypeOf((*MockBackendState)(nil).SecretBackendRotated), arg0, arg1)
}

// UpdateSecretBackend mocks base method.
func (m *MockBackendState) UpdateSecretBackend(arg0 state.UpdateSecretBackendParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecretBackend", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecretBackend indicates an expected call of UpdateSecretBackend.
func (mr *MockBackendStateMockRecorder) UpdateSecretBackend(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecretBackend", reflect.TypeOf((*MockBackendState)(nil).UpdateSecretBackend), arg0)
}
