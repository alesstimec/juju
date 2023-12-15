// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/cloud (interfaces: Backend,User,Model,ModelPoolBackend)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	credentialcommon "github.com/juju/juju/apiserver/common/credentialcommon"
	cloud "github.com/juju/juju/apiserver/facades/client/cloud"
	cloud0 "github.com/juju/juju/cloud"
	controller "github.com/juju/juju/controller"
	permission "github.com/juju/juju/core/permission"
	config "github.com/juju/juju/environs/config"
	context "github.com/juju/juju/environs/context"
	state "github.com/juju/juju/state"
	names "github.com/juju/names/v5"
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

// AddCloud mocks base method.
func (m *MockBackend) AddCloud(arg0 cloud0.Cloud, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCloud", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCloud indicates an expected call of AddCloud.
func (mr *MockBackendMockRecorder) AddCloud(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCloud", reflect.TypeOf((*MockBackend)(nil).AddCloud), arg0, arg1)
}

// AllCloudCredentials mocks base method.
func (m *MockBackend) AllCloudCredentials(arg0 names.UserTag) ([]state.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllCloudCredentials", arg0)
	ret0, _ := ret[0].([]state.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllCloudCredentials indicates an expected call of AllCloudCredentials.
func (mr *MockBackendMockRecorder) AllCloudCredentials(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllCloudCredentials", reflect.TypeOf((*MockBackend)(nil).AllCloudCredentials), arg0)
}

// Cloud mocks base method.
func (m *MockBackend) Cloud(arg0 string) (cloud0.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cloud", arg0)
	ret0, _ := ret[0].(cloud0.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cloud indicates an expected call of Cloud.
func (mr *MockBackendMockRecorder) Cloud(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cloud", reflect.TypeOf((*MockBackend)(nil).Cloud), arg0)
}

// CloudCredential mocks base method.
func (m *MockBackend) CloudCredential(arg0 names.CloudCredentialTag) (state.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential", arg0)
	ret0, _ := ret[0].(state.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockBackendMockRecorder) CloudCredential(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockBackend)(nil).CloudCredential), arg0)
}

// CloudCredentials mocks base method.
func (m *MockBackend) CloudCredentials(arg0 names.UserTag, arg1 string) (map[string]state.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredentials", arg0, arg1)
	ret0, _ := ret[0].(map[string]state.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredentials indicates an expected call of CloudCredentials.
func (mr *MockBackendMockRecorder) CloudCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredentials", reflect.TypeOf((*MockBackend)(nil).CloudCredentials), arg0, arg1)
}

// Clouds mocks base method.
func (m *MockBackend) Clouds() (map[names.CloudTag]cloud0.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clouds")
	ret0, _ := ret[0].(map[names.CloudTag]cloud0.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Clouds indicates an expected call of Clouds.
func (mr *MockBackendMockRecorder) Clouds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clouds", reflect.TypeOf((*MockBackend)(nil).Clouds))
}

// CloudsForUser mocks base method.
func (m *MockBackend) CloudsForUser(arg0 names.UserTag, arg1 bool) ([]state.CloudInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudsForUser", arg0, arg1)
	ret0, _ := ret[0].([]state.CloudInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudsForUser indicates an expected call of CloudsForUser.
func (mr *MockBackendMockRecorder) CloudsForUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudsForUser", reflect.TypeOf((*MockBackend)(nil).CloudsForUser), arg0, arg1)
}

// ControllerConfig mocks base method.
func (m *MockBackend) ControllerConfig() (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller.Config)
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

// CredentialModels mocks base method.
func (m *MockBackend) CredentialModels(arg0 names.CloudCredentialTag) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CredentialModels", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CredentialModels indicates an expected call of CredentialModels.
func (mr *MockBackendMockRecorder) CredentialModels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CredentialModels", reflect.TypeOf((*MockBackend)(nil).CredentialModels), arg0)
}

// CredentialModelsAndOwnerAccess mocks base method.
func (m *MockBackend) CredentialModelsAndOwnerAccess(arg0 names.CloudCredentialTag) ([]state.CredentialOwnerModelAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CredentialModelsAndOwnerAccess", arg0)
	ret0, _ := ret[0].([]state.CredentialOwnerModelAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CredentialModelsAndOwnerAccess indicates an expected call of CredentialModelsAndOwnerAccess.
func (mr *MockBackendMockRecorder) CredentialModelsAndOwnerAccess(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CredentialModelsAndOwnerAccess", reflect.TypeOf((*MockBackend)(nil).CredentialModelsAndOwnerAccess), arg0)
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

// GetCloudUsers mocks base method.
func (m *MockBackend) GetCloudUsers(arg0 string) (map[string]permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudUsers", arg0)
	ret0, _ := ret[0].(map[string]permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudUsers indicates an expected call of GetCloudUsers.
func (mr *MockBackendMockRecorder) GetCloudUsers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudUsers", reflect.TypeOf((*MockBackend)(nil).GetCloudUsers), arg0)
}

// Model mocks base method.
func (m *MockBackend) Model() (cloud.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(cloud.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Model indicates an expected call of Model.
func (mr *MockBackendMockRecorder) Model() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockBackend)(nil).Model))
}

// ModelConfig mocks base method.
func (m *MockBackend) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockBackendMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockBackend)(nil).ModelConfig))
}

// RemoveCloud mocks base method.
func (m *MockBackend) RemoveCloud(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCloud", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCloud indicates an expected call of RemoveCloud.
func (mr *MockBackendMockRecorder) RemoveCloud(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCloud", reflect.TypeOf((*MockBackend)(nil).RemoveCloud), arg0)
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

// RemoveCloudCredential mocks base method.
func (m *MockBackend) RemoveCloudCredential(arg0 names.CloudCredentialTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCloudCredential", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCloudCredential indicates an expected call of RemoveCloudCredential.
func (mr *MockBackendMockRecorder) RemoveCloudCredential(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCloudCredential", reflect.TypeOf((*MockBackend)(nil).RemoveCloudCredential), arg0)
}

// RemoveModelsCredential mocks base method.
func (m *MockBackend) RemoveModelsCredential(arg0 names.CloudCredentialTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveModelsCredential", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveModelsCredential indicates an expected call of RemoveModelsCredential.
func (mr *MockBackendMockRecorder) RemoveModelsCredential(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveModelsCredential", reflect.TypeOf((*MockBackend)(nil).RemoveModelsCredential), arg0)
}

// UpdateCloud mocks base method.
func (m *MockBackend) UpdateCloud(arg0 cloud0.Cloud) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCloud", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCloud indicates an expected call of UpdateCloud.
func (mr *MockBackendMockRecorder) UpdateCloud(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCloud", reflect.TypeOf((*MockBackend)(nil).UpdateCloud), arg0)
}

// UpdateCloudAccess mocks base method.
func (m *MockBackend) UpdateCloudAccess(arg0 string, arg1 names.UserTag, arg2 permission.Access) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCloudAccess", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCloudAccess indicates an expected call of UpdateCloudAccess.
func (mr *MockBackendMockRecorder) UpdateCloudAccess(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCloudAccess", reflect.TypeOf((*MockBackend)(nil).UpdateCloudAccess), arg0, arg1, arg2)
}

// UpdateCloudCredential mocks base method.
func (m *MockBackend) UpdateCloudCredential(arg0 names.CloudCredentialTag, arg1 cloud0.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCloudCredential", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCloudCredential indicates an expected call of UpdateCloudCredential.
func (mr *MockBackendMockRecorder) UpdateCloudCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCloudCredential", reflect.TypeOf((*MockBackend)(nil).UpdateCloudCredential), arg0, arg1)
}

// User mocks base method.
func (m *MockBackend) User(arg0 names.UserTag) (cloud.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User", arg0)
	ret0, _ := ret[0].(cloud.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// User indicates an expected call of User.
func (mr *MockBackendMockRecorder) User(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockBackend)(nil).User), arg0)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// DisplayName mocks base method.
func (m *MockUser) DisplayName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisplayName")
	ret0, _ := ret[0].(string)
	return ret0
}

// DisplayName indicates an expected call of DisplayName.
func (mr *MockUserMockRecorder) DisplayName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisplayName", reflect.TypeOf((*MockUser)(nil).DisplayName))
}

// MockModel is a mock of Model interface.
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel.
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance.
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// Cloud mocks base method.
func (m *MockModel) Cloud() (cloud0.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cloud")
	ret0, _ := ret[0].(cloud0.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cloud indicates an expected call of Cloud.
func (mr *MockModelMockRecorder) Cloud() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cloud", reflect.TypeOf((*MockModel)(nil).Cloud))
}

// CloudCredential mocks base method.
func (m *MockModel) CloudCredential() (state.Credential, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential")
	ret0, _ := ret[0].(state.Credential)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockModelMockRecorder) CloudCredential() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockModel)(nil).CloudCredential))
}

// CloudName mocks base method.
func (m *MockModel) CloudName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudName")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudName indicates an expected call of CloudName.
func (mr *MockModelMockRecorder) CloudName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudName", reflect.TypeOf((*MockModel)(nil).CloudName))
}

// CloudRegion mocks base method.
func (m *MockModel) CloudRegion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudRegion")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudRegion indicates an expected call of CloudRegion.
func (mr *MockModelMockRecorder) CloudRegion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudRegion", reflect.TypeOf((*MockModel)(nil).CloudRegion))
}

// UUID mocks base method.
func (m *MockModel) UUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UUID indicates an expected call of UUID.
func (mr *MockModelMockRecorder) UUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UUID", reflect.TypeOf((*MockModel)(nil).UUID))
}

// MockModelPoolBackend is a mock of ModelPoolBackend interface.
type MockModelPoolBackend struct {
	ctrl     *gomock.Controller
	recorder *MockModelPoolBackendMockRecorder
}

// MockModelPoolBackendMockRecorder is the mock recorder for MockModelPoolBackend.
type MockModelPoolBackendMockRecorder struct {
	mock *MockModelPoolBackend
}

// NewMockModelPoolBackend creates a new mock instance.
func NewMockModelPoolBackend(ctrl *gomock.Controller) *MockModelPoolBackend {
	mock := &MockModelPoolBackend{ctrl: ctrl}
	mock.recorder = &MockModelPoolBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelPoolBackend) EXPECT() *MockModelPoolBackendMockRecorder {
	return m.recorder
}

// GetModelCallContext mocks base method.
func (m *MockModelPoolBackend) GetModelCallContext(arg0 string) (credentialcommon.PersistentBackend, context.ProviderCallContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModelCallContext", arg0)
	ret0, _ := ret[0].(credentialcommon.PersistentBackend)
	ret1, _ := ret[1].(context.ProviderCallContext)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetModelCallContext indicates an expected call of GetModelCallContext.
func (mr *MockModelPoolBackendMockRecorder) GetModelCallContext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModelCallContext", reflect.TypeOf((*MockModelPoolBackend)(nil).GetModelCallContext), arg0)
}

// SystemState mocks base method.
func (m *MockModelPoolBackend) SystemState() (*state.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SystemState")
	ret0, _ := ret[0].(*state.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SystemState indicates an expected call of SystemState.
func (mr *MockModelPoolBackendMockRecorder) SystemState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SystemState", reflect.TypeOf((*MockModelPoolBackend)(nil).SystemState))
}
