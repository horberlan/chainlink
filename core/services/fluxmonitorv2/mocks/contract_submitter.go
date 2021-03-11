// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	mock "github.com/stretchr/testify/mock"
)

// ContractSubmitter is an autogenerated mock type for the ContractSubmitter type
type ContractSubmitter struct {
	mock.Mock
}

// Submit provides a mock function with given fields: roundID, submission
func (_m *ContractSubmitter) Submit(roundID *big.Int, submission *big.Int) error {
	ret := _m.Called(roundID, submission)

	var r0 error
	if rf, ok := ret.Get(0).(func(*big.Int, *big.Int) error); ok {
		r0 = rf(roundID, submission)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
