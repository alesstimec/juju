// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/controller (interfaces: Backend,Application,Relation)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	charm "github.com/juju/charm/v10"
	controller "github.com/juju/juju/apiserver/facades/client/controller"
	controller0 "github.com/juju/juju/controller"
	permission "github.com/juju/juju/core/permission"
	state "github.com/juju/juju/state"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// AddControllerUser mocks base method.
func (m *MockBackend) AddControllerUser(arg0 state.UserAccessSpec) (permission.UserAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddControllerUser", arg0)
	ret0, _ := ret[0].(permission.UserAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddControllerUser indicates an expected call of AddControllerUser.
func (mr *MockBackendMockRecorder) AddControllerUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddControllerUser", reflect.TypeOf((*MockBackend)(nil).AddControllerUser), arg0)
}

// AllBlocksForController mocks base method.
func (m *MockBackend) AllBlocksForController() ([]state.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllBlocksForController")
	ret0, _ := ret[0].([]state.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllBlocksForController indicates an expected call of AllBlocksForController.
func (mr *MockBackendMockRecorder) AllBlocksForController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllBlocksForController", reflect.TypeOf((*MockBackend)(nil).AllBlocksForController))
}

// AllModelUUIDs mocks base method.
func (m *MockBackend) AllModelUUIDs() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllModelUUIDs")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllModelUUIDs indicates an expected call of AllModelUUIDs.
func (mr *MockBackendMockRecorder) AllModelUUIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllModelUUIDs", reflect.TypeOf((*MockBackend)(nil).AllModelUUIDs))
}

// Application mocks base method.
func (m *MockBackend) Application(arg0 string) (controller.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application", arg0)
	ret0, _ := ret[0].(controller.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application.
func (mr *MockBackendMockRecorder) Application(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockBackend)(nil).Application), arg0)
}

// ControllerConfig mocks base method.
func (m *MockBackend) ControllerConfig() (controller0.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller0.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockBackendMockRecorder) ControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockBackend)(nil).ControllerConfig))
}

// ControllerInfo mocks base method.
func (m *MockBackend) ControllerInfo() (*state.ControllerInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerInfo")
	ret0, _ := ret[0].(*state.ControllerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerInfo indicates an expected call of ControllerInfo.
func (mr *MockBackendMockRecorder) ControllerInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerInfo", reflect.TypeOf((*MockBackend)(nil).ControllerInfo))
}

// ControllerModelUUID mocks base method.
func (m *MockBackend) ControllerModelUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerModelUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ControllerModelUUID indicates an expected call of ControllerModelUUID.
func (mr *MockBackendMockRecorder) ControllerModelUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerModelUUID", reflect.TypeOf((*MockBackend)(nil).ControllerModelUUID))
}

// ControllerTag mocks base method.
func (m *MockBackend) ControllerTag() names.ControllerTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerTag")
	ret0, _ := ret[0].(names.ControllerTag)
	return ret0
}

// ControllerTag indicates an expected call of ControllerTag.
func (mr *MockBackendMockRecorder) ControllerTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerTag", reflect.TypeOf((*MockBackend)(nil).ControllerTag))
}

// CreateCloudAccess mocks base method.
func (m *MockBackend) CreateCloudAccess(arg0 string, arg1 names.UserTag, arg2 permission.Access) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCloudAccess", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCloudAccess indicates an expected call of CreateCloudAccess.
func (mr *MockBackendMockRecorder) CreateCloudAccess(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCloudAccess", reflect.TypeOf((*MockBackend)(nil).CreateCloudAccess), arg0, arg1, arg2)
}

// GetCloudAccess mocks base method.
func (m *MockBackend) GetCloudAccess(arg0 string, arg1 names.UserTag) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudAccess", arg0, arg1)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudAccess indicates an expected call of GetCloudAccess.
func (mr *MockBackendMockRecorder) GetCloudAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudAccess", reflect.TypeOf((*MockBackend)(nil).GetCloudAccess), arg0, arg1)
}

// Model mocks base method.
func (m *MockBackend) Model() (*state.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(*state.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Model indicates an expected call of Model.
func (mr *MockBackendMockRecorder) Model() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockBackend)(nil).Model))
}

// ModelExists mocks base method.
func (m *MockBackend) ModelExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelExists indicates an expected call of ModelExists.
func (mr *MockBackendMockRecorder) ModelExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelExists", reflect.TypeOf((*MockBackend)(nil).ModelExists), arg0)
}

// MongoVersion mocks base method.
func (m *MockBackend) MongoVersion() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MongoVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MongoVersion indicates an expected call of MongoVersion.
func (mr *MockBackendMockRecorder) MongoVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MongoVersion", reflect.TypeOf((*MockBackend)(nil).MongoVersion))
}

// RemoveAllBlocksForController mocks base method.
func (m *MockBackend) RemoveAllBlocksForController() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAllBlocksForController")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAllBlocksForController indicates an expected call of RemoveAllBlocksForController.
func (mr *MockBackendMockRecorder) RemoveAllBlocksForController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAllBlocksForController", reflect.TypeOf((*MockBackend)(nil).RemoveAllBlocksForController))
}

// RemoveCloudAccess mocks base method.
func (m *MockBackend) RemoveCloudAccess(arg0 string, arg1 names.UserTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCloudAccess", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCloudAccess indicates an expected call of RemoveCloudAccess.
func (mr *MockBackendMockRecorder) RemoveCloudAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCloudAccess", reflect.TypeOf((*MockBackend)(nil).RemoveCloudAccess), arg0, arg1)
}

// RemoveUserAccess mocks base method.
func (m *MockBackend) RemoveUserAccess(arg0 names.UserTag, arg1 names.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserAccess", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserAccess indicates an expected call of RemoveUserAccess.
func (mr *MockBackendMockRecorder) RemoveUserAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserAccess", reflect.TypeOf((*MockBackend)(nil).RemoveUserAccess), arg0, arg1)
}

// SetUserAccess mocks base method.
func (m *MockBackend) SetUserAccess(arg0 names.UserTag, arg1 names.Tag, arg2 permission.Access) (permission.UserAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserAccess", arg0, arg1, arg2)
	ret0, _ := ret[0].(permission.UserAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetUserAccess indicates an expected call of SetUserAccess.
func (mr *MockBackendMockRecorder) SetUserAccess(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserAccess", reflect.TypeOf((*MockBackend)(nil).SetUserAccess), arg0, arg1, arg2)
}

// UpdateControllerConfig mocks base method.
func (m *MockBackend) UpdateControllerConfig(arg0 map[string]interface{}, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateControllerConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateControllerConfig indicates an expected call of UpdateControllerConfig.
func (mr *MockBackendMockRecorder) UpdateControllerConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateControllerConfig", reflect.TypeOf((*MockBackend)(nil).UpdateControllerConfig), arg0, arg1)
}

// UserAccess mocks base method.
func (m *MockBackend) UserAccess(arg0 names.UserTag, arg1 names.Tag) (permission.UserAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserAccess", arg0, arg1)
	ret0, _ := ret[0].(permission.UserAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserAccess indicates an expected call of UserAccess.
func (mr *MockBackendMockRecorder) UserAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAccess", reflect.TypeOf((*MockBackend)(nil).UserAccess), arg0, arg1)
}

// UserPermission mocks base method.
func (m *MockBackend) UserPermission(arg0 names.UserTag, arg1 names.Tag) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserPermission", arg0, arg1)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserPermission indicates an expected call of UserPermission.
func (mr *MockBackendMockRecorder) UserPermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserPermission", reflect.TypeOf((*MockBackend)(nil).UserPermission), arg0, arg1)
}

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// CharmConfig mocks base method.
func (m *MockApplication) CharmConfig(arg0 string) (charm.Settings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmConfig", arg0)
	ret0, _ := ret[0].(charm.Settings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CharmConfig indicates an expected call of CharmConfig.
func (mr *MockApplicationMockRecorder) CharmConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmConfig", reflect.TypeOf((*MockApplication)(nil).CharmConfig), arg0)
}

// Name mocks base method.
func (m *MockApplication) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockApplicationMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockApplication)(nil).Name))
}

// Relations mocks base method.
func (m *MockApplication) Relations() ([]controller.Relation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Relations")
	ret0, _ := ret[0].([]controller.Relation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Relations indicates an expected call of Relations.
func (mr *MockApplicationMockRecorder) Relations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Relations", reflect.TypeOf((*MockApplication)(nil).Relations))
}

// MockRelation is a mock of Relation interface.
type MockRelation struct {
	ctrl     *gomock.Controller
	recorder *MockRelationMockRecorder
}

// MockRelationMockRecorder is the mock recorder for MockRelation.
type MockRelationMockRecorder struct {
	mock *MockRelation
}

// NewMockRelation creates a new mock instance.
func NewMockRelation(ctrl *gomock.Controller) *MockRelation {
	mock := &MockRelation{ctrl: ctrl}
	mock.recorder = &MockRelationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRelation) EXPECT() *MockRelationMockRecorder {
	return m.recorder
}

// ApplicationSettings mocks base method.
func (m *MockRelation) ApplicationSettings(arg0 string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationSettings", arg0)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationSettings indicates an expected call of ApplicationSettings.
func (mr *MockRelationMockRecorder) ApplicationSettings(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationSettings", reflect.TypeOf((*MockRelation)(nil).ApplicationSettings), arg0)
}

// Endpoint mocks base method.
func (m *MockRelation) Endpoint(arg0 string) (state.Endpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Endpoint", arg0)
	ret0, _ := ret[0].(state.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Endpoint indicates an expected call of Endpoint.
func (mr *MockRelationMockRecorder) Endpoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Endpoint", reflect.TypeOf((*MockRelation)(nil).Endpoint), arg0)
}

// ModelUUID mocks base method.
func (m *MockRelation) ModelUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ModelUUID indicates an expected call of ModelUUID.
func (mr *MockRelationMockRecorder) ModelUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelUUID", reflect.TypeOf((*MockRelation)(nil).ModelUUID))
}

// RelatedEndpoints mocks base method.
func (m *MockRelation) RelatedEndpoints(arg0 string) ([]state.Endpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelatedEndpoints", arg0)
	ret0, _ := ret[0].([]state.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RelatedEndpoints indicates an expected call of RelatedEndpoints.
func (mr *MockRelationMockRecorder) RelatedEndpoints(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelatedEndpoints", reflect.TypeOf((*MockRelation)(nil).RelatedEndpoints), arg0)
}
