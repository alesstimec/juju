// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/secretbackendrotate (interfaces: SecretBackendManagerFacade)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	watcher "github.com/juju/juju/core/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockSecretBackendManagerFacade is a mock of SecretBackendManagerFacade interface.
type MockSecretBackendManagerFacade struct {
	ctrl     *gomock.Controller
	recorder *MockSecretBackendManagerFacadeMockRecorder
}

// MockSecretBackendManagerFacadeMockRecorder is the mock recorder for MockSecretBackendManagerFacade.
type MockSecretBackendManagerFacadeMockRecorder struct {
	mock *MockSecretBackendManagerFacade
}

// NewMockSecretBackendManagerFacade creates a new mock instance.
func NewMockSecretBackendManagerFacade(ctrl *gomock.Controller) *MockSecretBackendManagerFacade {
	mock := &MockSecretBackendManagerFacade{ctrl: ctrl}
	mock.recorder = &MockSecretBackendManagerFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretBackendManagerFacade) EXPECT() *MockSecretBackendManagerFacadeMockRecorder {
	return m.recorder
}

// RotateBackendTokens mocks base method.
func (m *MockSecretBackendManagerFacade) RotateBackendTokens(arg0 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RotateBackendTokens", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RotateBackendTokens indicates an expected call of RotateBackendTokens.
func (mr *MockSecretBackendManagerFacadeMockRecorder) RotateBackendTokens(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RotateBackendTokens", reflect.TypeOf((*MockSecretBackendManagerFacade)(nil).RotateBackendTokens), arg0...)
}

// WatchTokenRotationChanges mocks base method.
func (m *MockSecretBackendManagerFacade) WatchTokenRotationChanges() (watcher.SecretBackendRotateWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchTokenRotationChanges")
	ret0, _ := ret[0].(watcher.SecretBackendRotateWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchTokenRotationChanges indicates an expected call of WatchTokenRotationChanges.
func (mr *MockSecretBackendManagerFacadeMockRecorder) WatchTokenRotationChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchTokenRotationChanges", reflect.TypeOf((*MockSecretBackendManagerFacade)(nil).WatchTokenRotationChanges))
}
