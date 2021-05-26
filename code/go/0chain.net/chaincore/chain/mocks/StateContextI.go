// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	block "0chain.net/chaincore/block"

	encryption "0chain.net/core/encryption"

	mock "github.com/stretchr/testify/mock"

	state "0chain.net/chaincore/state"

	transaction "0chain.net/chaincore/transaction"

	util "0chain.net/core/util"
)

// StateContextI is an autogenerated mock type for the StateContextI type
type StateContextI struct {
	mock.Mock
}

// AddMint provides a mock function with given fields: m
func (_m *StateContextI) AddMint(m *state.Mint) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(*state.Mint) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddSignedTransfer provides a mock function with given fields: st
func (_m *StateContextI) AddSignedTransfer(st *state.SignedTransfer) {
	_m.Called(st)
}

// AddTransfer provides a mock function with given fields: t
func (_m *StateContextI) AddTransfer(t *state.Transfer) error {
	ret := _m.Called(t)

	var r0 error
	if rf, ok := ret.Get(0).(func(*state.Transfer) error); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTrieNode provides a mock function with given fields: key
func (_m *StateContextI) DeleteTrieNode(key string) (string, error) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlock provides a mock function with given fields:
func (_m *StateContextI) GetBlock() *block.Block {
	ret := _m.Called()

	var r0 *block.Block
	if rf, ok := ret.Get(0).(func() *block.Block); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*block.Block)
		}
	}

	return r0
}

// GetBlockSharders provides a mock function with given fields: b
func (_m *StateContextI) GetBlockSharders(b *block.Block) []string {
	ret := _m.Called(b)

	var r0 []string
	if rf, ok := ret.Get(0).(func(*block.Block) []string); ok {
		r0 = rf(b)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetChainCurrentMagicBlock provides a mock function with given fields:
func (_m *StateContextI) GetChainCurrentMagicBlock() *block.MagicBlock {
	ret := _m.Called()

	var r0 *block.MagicBlock
	if rf, ok := ret.Get(0).(func() *block.MagicBlock); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*block.MagicBlock)
		}
	}

	return r0
}

// GetClientBalance provides a mock function with given fields: clientID
func (_m *StateContextI) GetClientBalance(clientID string) (state.Balance, error) {
	ret := _m.Called(clientID)

	var r0 state.Balance
	if rf, ok := ret.Get(0).(func(string) state.Balance); ok {
		r0 = rf(clientID)
	} else {
		r0 = ret.Get(0).(state.Balance)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastestFinalizedMagicBlock provides a mock function with given fields:
func (_m *StateContextI) GetLastestFinalizedMagicBlock() *block.Block {
	ret := _m.Called()

	var r0 *block.Block
	if rf, ok := ret.Get(0).(func() *block.Block); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*block.Block)
		}
	}

	return r0
}

// GetMints provides a mock function with given fields:
func (_m *StateContextI) GetMints() []*state.Mint {
	ret := _m.Called()

	var r0 []*state.Mint
	if rf, ok := ret.Get(0).(func() []*state.Mint); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*state.Mint)
		}
	}

	return r0
}

// GetSignatureScheme provides a mock function with given fields:
func (_m *StateContextI) GetSignatureScheme() encryption.SignatureScheme {
	ret := _m.Called()

	var r0 encryption.SignatureScheme
	if rf, ok := ret.Get(0).(func() encryption.SignatureScheme); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(encryption.SignatureScheme)
		}
	}

	return r0
}

// GetSignedTransfers provides a mock function with given fields:
func (_m *StateContextI) GetSignedTransfers() []*state.SignedTransfer {
	ret := _m.Called()

	var r0 []*state.SignedTransfer
	if rf, ok := ret.Get(0).(func() []*state.SignedTransfer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*state.SignedTransfer)
		}
	}

	return r0
}

// GetState provides a mock function with given fields:
func (_m *StateContextI) GetState() util.MerklePatriciaTrieI {
	ret := _m.Called()

	var r0 util.MerklePatriciaTrieI
	if rf, ok := ret.Get(0).(func() util.MerklePatriciaTrieI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(util.MerklePatriciaTrieI)
		}
	}

	return r0
}

// GetTransaction provides a mock function with given fields:
func (_m *StateContextI) GetTransaction() *transaction.Transaction {
	ret := _m.Called()

	var r0 *transaction.Transaction
	if rf, ok := ret.Get(0).(func() *transaction.Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Transaction)
		}
	}

	return r0
}

// GetTransfers provides a mock function with given fields:
func (_m *StateContextI) GetTransfers() []*state.Transfer {
	ret := _m.Called()

	var r0 []*state.Transfer
	if rf, ok := ret.Get(0).(func() []*state.Transfer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*state.Transfer)
		}
	}

	return r0
}

// GetTrieNode provides a mock function with given fields: key
func (_m *StateContextI) GetTrieNode(key string) (util.Serializable, error) {
	ret := _m.Called(key)

	var r0 util.Serializable
	if rf, ok := ret.Get(0).(func(string) util.Serializable); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(util.Serializable)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertTrieNode provides a mock function with given fields: key, node
func (_m *StateContextI) InsertTrieNode(key string, node util.Serializable) (string, error) {
	ret := _m.Called(key, node)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, util.Serializable) string); ok {
		r0 = rf(key, node)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, util.Serializable) error); ok {
		r1 = rf(key, node)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetMagicBlock provides a mock function with given fields: _a0
func (_m *StateContextI) SetMagicBlock(_a0 *block.MagicBlock) {
	_m.Called(_a0)
}

// SetStateContext provides a mock function with given fields: st
func (_m *StateContextI) SetStateContext(st *state.State) error {
	ret := _m.Called(st)

	var r0 error
	if rf, ok := ret.Get(0).(func(*state.State) error); ok {
		r0 = rf(st)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Validate provides a mock function with given fields:
func (_m *StateContextI) Validate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
