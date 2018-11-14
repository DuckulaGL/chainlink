// Code generated by MockGen. DO NOT EDIT.
// Source: store/tx_manager.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	go_ethereum "github.com/ethereum/go-ethereum"
	accounts "github.com/ethereum/go-ethereum/accounts"
	common "github.com/ethereum/go-ethereum/common"
	gomock "github.com/golang/mock/gomock"
	store "github.com/smartcontractkit/chainlink/store"
	assets "github.com/smartcontractkit/chainlink/store/assets"
	models "github.com/smartcontractkit/chainlink/store/models"
	reflect "reflect"
)

// MockTxManager is a mock of TxManager interface
type MockTxManager struct {
	ctrl     *gomock.Controller
	recorder *MockTxManagerMockRecorder
}

// MockTxManagerMockRecorder is the mock recorder for MockTxManager
type MockTxManagerMockRecorder struct {
	mock *MockTxManager
}

// NewMockTxManager creates a new mock instance
func NewMockTxManager(ctrl *gomock.Controller) *MockTxManager {
	mock := &MockTxManager{ctrl: ctrl}
	mock.recorder = &MockTxManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTxManager) EXPECT() *MockTxManagerMockRecorder {
	return m.recorder
}

// CreateTx mocks base method
func (m *MockTxManager) CreateTx(to common.Address, data []byte) (*models.Tx, error) {
	ret := m.ctrl.Call(m, "CreateTx", to, data)
	ret0, _ := ret[0].(*models.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTx indicates an expected call of CreateTx
func (mr *MockTxManagerMockRecorder) CreateTx(to, data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockTxManager)(nil).CreateTx), to, data)
}

// ActivateAccount mocks base method
func (m *MockTxManager) ActivateAccount(account accounts.Account) error {
	ret := m.ctrl.Call(m, "ActivateAccount", account)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateAccount indicates an expected call of ActivateAccount
func (mr *MockTxManagerMockRecorder) ActivateAccount(account interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateAccount", reflect.TypeOf((*MockTxManager)(nil).ActivateAccount), account)
}

// MeetsMinConfirmations mocks base method
func (m *MockTxManager) MeetsMinConfirmations(hash common.Hash) (bool, error) {
	ret := m.ctrl.Call(m, "MeetsMinConfirmations", hash)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MeetsMinConfirmations indicates an expected call of MeetsMinConfirmations
func (mr *MockTxManagerMockRecorder) MeetsMinConfirmations(hash interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MeetsMinConfirmations", reflect.TypeOf((*MockTxManager)(nil).MeetsMinConfirmations), hash)
}

// WithdrawLink mocks base method
func (m *MockTxManager) WithdrawLink(wr models.WithdrawalRequest) (common.Hash, error) {
	ret := m.ctrl.Call(m, "WithdrawLink", wr)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WithdrawLink indicates an expected call of WithdrawLink
func (mr *MockTxManagerMockRecorder) WithdrawLink(wr interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawLink", reflect.TypeOf((*MockTxManager)(nil).WithdrawLink), wr)
}

// GetLinkBalance mocks base method
func (m *MockTxManager) GetLinkBalance(address common.Address) (*assets.Link, error) {
	ret := m.ctrl.Call(m, "GetLinkBalance", address)
	ret0, _ := ret[0].(*assets.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkBalance indicates an expected call of GetLinkBalance
func (mr *MockTxManagerMockRecorder) GetLinkBalance(address interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkBalance", reflect.TypeOf((*MockTxManager)(nil).GetLinkBalance), address)
}

// GetActiveAccount mocks base method
func (m *MockTxManager) GetActiveAccount() *store.ActiveAccount {
	ret := m.ctrl.Call(m, "GetActiveAccount")
	ret0, _ := ret[0].(*store.ActiveAccount)
	return ret0
}

// GetActiveAccount indicates an expected call of GetActiveAccount
func (mr *MockTxManagerMockRecorder) GetActiveAccount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveAccount", reflect.TypeOf((*MockTxManager)(nil).GetActiveAccount))
}

// GetEthBalance mocks base method
func (m *MockTxManager) GetEthBalance(address common.Address) (*assets.Eth, error) {
	ret := m.ctrl.Call(m, "GetEthBalance", address)
	ret0, _ := ret[0].(*assets.Eth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEthBalance indicates an expected call of GetEthBalance
func (mr *MockTxManagerMockRecorder) GetEthBalance(address interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEthBalance", reflect.TypeOf((*MockTxManager)(nil).GetEthBalance), address)
}

// SubscribeToNewHeads mocks base method
func (m *MockTxManager) SubscribeToNewHeads(channel chan<- models.BlockHeader) (models.EthSubscription, error) {
	ret := m.ctrl.Call(m, "SubscribeToNewHeads", channel)
	ret0, _ := ret[0].(models.EthSubscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToNewHeads indicates an expected call of SubscribeToNewHeads
func (mr *MockTxManagerMockRecorder) SubscribeToNewHeads(channel interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToNewHeads", reflect.TypeOf((*MockTxManager)(nil).SubscribeToNewHeads), channel)
}

// GetBlockByNumber mocks base method
func (m *MockTxManager) GetBlockByNumber(hex string) (models.BlockHeader, error) {
	ret := m.ctrl.Call(m, "GetBlockByNumber", hex)
	ret0, _ := ret[0].(models.BlockHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber
func (mr *MockTxManagerMockRecorder) GetBlockByNumber(hex interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockTxManager)(nil).GetBlockByNumber), hex)
}

// SubscribeToLogs mocks base method
func (m *MockTxManager) SubscribeToLogs(channel chan<- store.Log, q go_ethereum.FilterQuery) (models.EthSubscription, error) {
	ret := m.ctrl.Call(m, "SubscribeToLogs", channel, q)
	ret0, _ := ret[0].(models.EthSubscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToLogs indicates an expected call of SubscribeToLogs
func (mr *MockTxManagerMockRecorder) SubscribeToLogs(channel, q interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToLogs", reflect.TypeOf((*MockTxManager)(nil).SubscribeToLogs), channel, q)
}

// GetLogs mocks base method
func (m *MockTxManager) GetLogs(q go_ethereum.FilterQuery) ([]store.Log, error) {
	ret := m.ctrl.Call(m, "GetLogs", q)
	ret0, _ := ret[0].([]store.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs
func (mr *MockTxManagerMockRecorder) GetLogs(q interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockTxManager)(nil).GetLogs), q)
}
