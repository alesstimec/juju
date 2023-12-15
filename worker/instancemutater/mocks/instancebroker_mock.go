// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/instancemutater (interfaces: InstanceMutaterAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	instancemutater "github.com/juju/juju/api/agent/instancemutater"
	watcher "github.com/juju/juju/core/watcher"
	names "github.com/juju/names/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockInstanceMutaterAPI is a mock of InstanceMutaterAPI interface.
type MockInstanceMutaterAPI struct {
	ctrl     *gomock.Controller
	recorder *MockInstanceMutaterAPIMockRecorder
}

// MockInstanceMutaterAPIMockRecorder is the mock recorder for MockInstanceMutaterAPI.
type MockInstanceMutaterAPIMockRecorder struct {
	mock *MockInstanceMutaterAPI
}

// NewMockInstanceMutaterAPI creates a new mock instance.
func NewMockInstanceMutaterAPI(ctrl *gomock.Controller) *MockInstanceMutaterAPI {
	mock := &MockInstanceMutaterAPI{ctrl: ctrl}
	mock.recorder = &MockInstanceMutaterAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstanceMutaterAPI) EXPECT() *MockInstanceMutaterAPIMockRecorder {
	return m.recorder
}

// Machine mocks base method.
func (m *MockInstanceMutaterAPI) Machine(arg0 names.MachineTag) (instancemutater.MutaterMachine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Machine", arg0)
	ret0, _ := ret[0].(instancemutater.MutaterMachine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Machine indicates an expected call of Machine.
func (mr *MockInstanceMutaterAPIMockRecorder) Machine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machine", reflect.TypeOf((*MockInstanceMutaterAPI)(nil).Machine), arg0)
}

// WatchModelMachines mocks base method.
func (m *MockInstanceMutaterAPI) WatchModelMachines() (watcher.StringsWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchModelMachines")
	ret0, _ := ret[0].(watcher.StringsWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchModelMachines indicates an expected call of WatchModelMachines.
func (mr *MockInstanceMutaterAPIMockRecorder) WatchModelMachines() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchModelMachines", reflect.TypeOf((*MockInstanceMutaterAPI)(nil).WatchModelMachines))
}
