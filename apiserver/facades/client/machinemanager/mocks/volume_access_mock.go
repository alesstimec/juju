// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common/storagecommon (interfaces: VolumeAccess)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	state "github.com/juju/juju/state"
	names "github.com/juju/names/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockVolumeAccess is a mock of VolumeAccess interface.
type MockVolumeAccess struct {
	ctrl     *gomock.Controller
	recorder *MockVolumeAccessMockRecorder
}

// MockVolumeAccessMockRecorder is the mock recorder for MockVolumeAccess.
type MockVolumeAccessMockRecorder struct {
	mock *MockVolumeAccess
}

// NewMockVolumeAccess creates a new mock instance.
func NewMockVolumeAccess(ctrl *gomock.Controller) *MockVolumeAccess {
	mock := &MockVolumeAccess{ctrl: ctrl}
	mock.recorder = &MockVolumeAccessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVolumeAccess) EXPECT() *MockVolumeAccessMockRecorder {
	return m.recorder
}

// BlockDevices mocks base method.
func (m *MockVolumeAccess) BlockDevices(arg0 names.MachineTag) ([]state.BlockDeviceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockDevices", arg0)
	ret0, _ := ret[0].([]state.BlockDeviceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockDevices indicates an expected call of BlockDevices.
func (mr *MockVolumeAccessMockRecorder) BlockDevices(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockDevices", reflect.TypeOf((*MockVolumeAccess)(nil).BlockDevices), arg0)
}

// StorageInstanceVolume mocks base method.
func (m *MockVolumeAccess) StorageInstanceVolume(arg0 names.StorageTag) (state.Volume, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageInstanceVolume", arg0)
	ret0, _ := ret[0].(state.Volume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageInstanceVolume indicates an expected call of StorageInstanceVolume.
func (mr *MockVolumeAccessMockRecorder) StorageInstanceVolume(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageInstanceVolume", reflect.TypeOf((*MockVolumeAccess)(nil).StorageInstanceVolume), arg0)
}

// VolumeAttachment mocks base method.
func (m *MockVolumeAccess) VolumeAttachment(arg0 names.Tag, arg1 names.VolumeTag) (state.VolumeAttachment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeAttachment", arg0, arg1)
	ret0, _ := ret[0].(state.VolumeAttachment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VolumeAttachment indicates an expected call of VolumeAttachment.
func (mr *MockVolumeAccessMockRecorder) VolumeAttachment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeAttachment", reflect.TypeOf((*MockVolumeAccess)(nil).VolumeAttachment), arg0, arg1)
}

// VolumeAttachmentPlan mocks base method.
func (m *MockVolumeAccess) VolumeAttachmentPlan(arg0 names.Tag, arg1 names.VolumeTag) (state.VolumeAttachmentPlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeAttachmentPlan", arg0, arg1)
	ret0, _ := ret[0].(state.VolumeAttachmentPlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VolumeAttachmentPlan indicates an expected call of VolumeAttachmentPlan.
func (mr *MockVolumeAccessMockRecorder) VolumeAttachmentPlan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeAttachmentPlan", reflect.TypeOf((*MockVolumeAccess)(nil).VolumeAttachmentPlan), arg0, arg1)
}
