// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common/secrets (interfaces: Model,Credential,SecretsConsumer,SecretsMetaState,SecretsRemoveState)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	secrets "github.com/juju/juju/apiserver/common/secrets"
	cloud "github.com/juju/juju/cloud"
	secrets0 "github.com/juju/juju/core/secrets"
	config "github.com/juju/juju/environs/config"
	state "github.com/juju/juju/state"
	names "github.com/juju/names/v5"
	gomock "go.uber.org/mock/gomock"
)

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
func (m *MockModel) Cloud() (cloud.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cloud")
	ret0, _ := ret[0].(cloud.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cloud indicates an expected call of Cloud.
func (mr *MockModelMockRecorder) Cloud() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cloud", reflect.TypeOf((*MockModel)(nil).Cloud))
}

// CloudCredential mocks base method.
func (m *MockModel) CloudCredential() (secrets.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential")
	ret0, _ := ret[0].(secrets.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockModelMockRecorder) CloudCredential() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockModel)(nil).CloudCredential))
}

// Config mocks base method.
func (m *MockModel) Config() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Config indicates an expected call of Config.
func (mr *MockModelMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockModel)(nil).Config))
}

// ControllerUUID mocks base method.
func (m *MockModel) ControllerUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ControllerUUID indicates an expected call of ControllerUUID.
func (mr *MockModelMockRecorder) ControllerUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerUUID", reflect.TypeOf((*MockModel)(nil).ControllerUUID))
}

// ModelConfig mocks base method.
func (m *MockModel) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockModelMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockModel)(nil).ModelConfig))
}

// Name mocks base method.
func (m *MockModel) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockModelMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockModel)(nil).Name))
}

// State mocks base method.
func (m *MockModel) State() *state.State {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(*state.State)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockModelMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockModel)(nil).State))
}

// Type mocks base method.
func (m *MockModel) Type() state.ModelType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(state.ModelType)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockModelMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockModel)(nil).Type))
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

// WatchForModelConfigChanges mocks base method.
func (m *MockModel) WatchForModelConfigChanges() state.NotifyWatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchForModelConfigChanges")
	ret0, _ := ret[0].(state.NotifyWatcher)
	return ret0
}

// WatchForModelConfigChanges indicates an expected call of WatchForModelConfigChanges.
func (mr *MockModelMockRecorder) WatchForModelConfigChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchForModelConfigChanges", reflect.TypeOf((*MockModel)(nil).WatchForModelConfigChanges))
}

// MockCredential is a mock of Credential interface.
type MockCredential struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialMockRecorder
}

// MockCredentialMockRecorder is the mock recorder for MockCredential.
type MockCredentialMockRecorder struct {
	mock *MockCredential
}

// NewMockCredential creates a new mock instance.
func NewMockCredential(ctrl *gomock.Controller) *MockCredential {
	mock := &MockCredential{ctrl: ctrl}
	mock.recorder = &MockCredentialMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredential) EXPECT() *MockCredentialMockRecorder {
	return m.recorder
}

// Attributes mocks base method.
func (m *MockCredential) Attributes() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// Attributes indicates an expected call of Attributes.
func (mr *MockCredentialMockRecorder) Attributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockCredential)(nil).Attributes))
}

// AuthType mocks base method.
func (m *MockCredential) AuthType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthType")
	ret0, _ := ret[0].(string)
	return ret0
}

// AuthType indicates an expected call of AuthType.
func (mr *MockCredentialMockRecorder) AuthType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthType", reflect.TypeOf((*MockCredential)(nil).AuthType))
}

// MockSecretsConsumer is a mock of SecretsConsumer interface.
type MockSecretsConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsConsumerMockRecorder
}

// MockSecretsConsumerMockRecorder is the mock recorder for MockSecretsConsumer.
type MockSecretsConsumerMockRecorder struct {
	mock *MockSecretsConsumer
}

// NewMockSecretsConsumer creates a new mock instance.
func NewMockSecretsConsumer(ctrl *gomock.Controller) *MockSecretsConsumer {
	mock := &MockSecretsConsumer{ctrl: ctrl}
	mock.recorder = &MockSecretsConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsConsumer) EXPECT() *MockSecretsConsumerMockRecorder {
	return m.recorder
}

// SecretAccess mocks base method.
func (m *MockSecretsConsumer) SecretAccess(arg0 *secrets0.URI, arg1 names.Tag) (secrets0.SecretRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretAccess", arg0, arg1)
	ret0, _ := ret[0].(secrets0.SecretRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SecretAccess indicates an expected call of SecretAccess.
func (mr *MockSecretsConsumerMockRecorder) SecretAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretAccess", reflect.TypeOf((*MockSecretsConsumer)(nil).SecretAccess), arg0, arg1)
}

// MockSecretsMetaState is a mock of SecretsMetaState interface.
type MockSecretsMetaState struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsMetaStateMockRecorder
}

// MockSecretsMetaStateMockRecorder is the mock recorder for MockSecretsMetaState.
type MockSecretsMetaStateMockRecorder struct {
	mock *MockSecretsMetaState
}

// NewMockSecretsMetaState creates a new mock instance.
func NewMockSecretsMetaState(ctrl *gomock.Controller) *MockSecretsMetaState {
	mock := &MockSecretsMetaState{ctrl: ctrl}
	mock.recorder = &MockSecretsMetaStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsMetaState) EXPECT() *MockSecretsMetaStateMockRecorder {
	return m.recorder
}

// ChangeSecretBackend mocks base method.
func (m *MockSecretsMetaState) ChangeSecretBackend(arg0 state.ChangeSecretBackendParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeSecretBackend", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeSecretBackend indicates an expected call of ChangeSecretBackend.
func (mr *MockSecretsMetaStateMockRecorder) ChangeSecretBackend(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeSecretBackend", reflect.TypeOf((*MockSecretsMetaState)(nil).ChangeSecretBackend), arg0)
}

// ListSecretRevisions mocks base method.
func (m *MockSecretsMetaState) ListSecretRevisions(arg0 *secrets0.URI) ([]*secrets0.SecretRevisionMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecretRevisions", arg0)
	ret0, _ := ret[0].([]*secrets0.SecretRevisionMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecretRevisions indicates an expected call of ListSecretRevisions.
func (mr *MockSecretsMetaStateMockRecorder) ListSecretRevisions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecretRevisions", reflect.TypeOf((*MockSecretsMetaState)(nil).ListSecretRevisions), arg0)
}

// ListSecrets mocks base method.
func (m *MockSecretsMetaState) ListSecrets(arg0 state.SecretsFilter) ([]*secrets0.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", arg0)
	ret0, _ := ret[0].([]*secrets0.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets.
func (mr *MockSecretsMetaStateMockRecorder) ListSecrets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecretsMetaState)(nil).ListSecrets), arg0)
}

// SecretGrants mocks base method.
func (m *MockSecretsMetaState) SecretGrants(arg0 *secrets0.URI, arg1 secrets0.SecretRole) ([]secrets0.AccessInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretGrants", arg0, arg1)
	ret0, _ := ret[0].([]secrets0.AccessInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SecretGrants indicates an expected call of SecretGrants.
func (mr *MockSecretsMetaStateMockRecorder) SecretGrants(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretGrants", reflect.TypeOf((*MockSecretsMetaState)(nil).SecretGrants), arg0, arg1)
}

// MockSecretsRemoveState is a mock of SecretsRemoveState interface.
type MockSecretsRemoveState struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsRemoveStateMockRecorder
}

// MockSecretsRemoveStateMockRecorder is the mock recorder for MockSecretsRemoveState.
type MockSecretsRemoveStateMockRecorder struct {
	mock *MockSecretsRemoveState
}

// NewMockSecretsRemoveState creates a new mock instance.
func NewMockSecretsRemoveState(ctrl *gomock.Controller) *MockSecretsRemoveState {
	mock := &MockSecretsRemoveState{ctrl: ctrl}
	mock.recorder = &MockSecretsRemoveStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsRemoveState) EXPECT() *MockSecretsRemoveStateMockRecorder {
	return m.recorder
}

// DeleteSecret mocks base method.
func (m *MockSecretsRemoveState) DeleteSecret(arg0 *secrets0.URI, arg1 ...int) ([]secrets0.ValueRef, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSecret", varargs...)
	ret0, _ := ret[0].([]secrets0.ValueRef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretsRemoveStateMockRecorder) DeleteSecret(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretsRemoveState)(nil).DeleteSecret), varargs...)
}

// GetSecret mocks base method.
func (m *MockSecretsRemoveState) GetSecret(arg0 *secrets0.URI) (*secrets0.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", arg0)
	ret0, _ := ret[0].(*secrets0.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretsRemoveStateMockRecorder) GetSecret(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretsRemoveState)(nil).GetSecret), arg0)
}

// GetSecretRevision mocks base method.
func (m *MockSecretsRemoveState) GetSecretRevision(arg0 *secrets0.URI, arg1 int) (*secrets0.SecretRevisionMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretRevision", arg0, arg1)
	ret0, _ := ret[0].(*secrets0.SecretRevisionMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretRevision indicates an expected call of GetSecretRevision.
func (mr *MockSecretsRemoveStateMockRecorder) GetSecretRevision(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretRevision", reflect.TypeOf((*MockSecretsRemoveState)(nil).GetSecretRevision), arg0, arg1)
}

// ListSecretRevisions mocks base method.
func (m *MockSecretsRemoveState) ListSecretRevisions(arg0 *secrets0.URI) ([]*secrets0.SecretRevisionMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecretRevisions", arg0)
	ret0, _ := ret[0].([]*secrets0.SecretRevisionMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecretRevisions indicates an expected call of ListSecretRevisions.
func (mr *MockSecretsRemoveStateMockRecorder) ListSecretRevisions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecretRevisions", reflect.TypeOf((*MockSecretsRemoveState)(nil).ListSecretRevisions), arg0)
}

// ListSecrets mocks base method.
func (m *MockSecretsRemoveState) ListSecrets(arg0 state.SecretsFilter) ([]*secrets0.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", arg0)
	ret0, _ := ret[0].([]*secrets0.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets.
func (mr *MockSecretsRemoveStateMockRecorder) ListSecrets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecretsRemoveState)(nil).ListSecrets), arg0)
}
