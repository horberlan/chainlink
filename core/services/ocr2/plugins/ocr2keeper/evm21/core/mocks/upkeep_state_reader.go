// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

// UpkeepStateReader is an autogenerated mock type for the UpkeepStateReader type
type UpkeepStateReader struct {
	mock.Mock
}

// SelectByWorkIDsInRange provides a mock function with given fields: ctx, start, end, workIDs
func (_m *UpkeepStateReader) SelectByWorkIDsInRange(ctx context.Context, start int64, end int64, workIDs ...string) ([]types.UpkeepState, error) {
	_va := make([]interface{}, len(workIDs))
	for _i := range workIDs {
		_va[_i] = workIDs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, start, end)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []types.UpkeepState
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, ...string) ([]types.UpkeepState, error)); ok {
		return rf(ctx, start, end, workIDs...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, ...string) []types.UpkeepState); ok {
		r0 = rf(ctx, start, end, workIDs...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.UpkeepState)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, ...string) error); ok {
		r1 = rf(ctx, start, end, workIDs...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUpkeepStateReader interface {
	mock.TestingT
	Cleanup(func())
}

// NewUpkeepStateReader creates a new instance of UpkeepStateReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUpkeepStateReader(t mockConstructorTestingTNewUpkeepStateReader) *UpkeepStateReader {
	mock := &UpkeepStateReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}