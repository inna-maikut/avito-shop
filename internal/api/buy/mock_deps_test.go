// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package buy -typed -destination mock_deps_test.go
//

// Package buy is a generated GoMock package.
package buy

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockbuying is a mock of buying interface.
type Mockbuying struct {
	ctrl     *gomock.Controller
	recorder *MockbuyingMockRecorder
}

// MockbuyingMockRecorder is the mock recorder for Mockbuying.
type MockbuyingMockRecorder struct {
	mock *Mockbuying
}

// NewMockbuying creates a new mock instance.
func NewMockbuying(ctrl *gomock.Controller) *Mockbuying {
	mock := &Mockbuying{ctrl: ctrl}
	mock.recorder = &MockbuyingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockbuying) EXPECT() *MockbuyingMockRecorder {
	return m.recorder
}

// Buy mocks base method.
func (m *Mockbuying) Buy(ctx context.Context, employeeID int64, merchName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Buy", ctx, employeeID, merchName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Buy indicates an expected call of Buy.
func (mr *MockbuyingMockRecorder) Buy(ctx, employeeID, merchName any) *MockbuyingBuyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buy", reflect.TypeOf((*Mockbuying)(nil).Buy), ctx, employeeID, merchName)
	return &MockbuyingBuyCall{Call: call}
}

// MockbuyingBuyCall wrap *gomock.Call
type MockbuyingBuyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockbuyingBuyCall) Return(arg0 error) *MockbuyingBuyCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockbuyingBuyCall) Do(f func(context.Context, int64, string) error) *MockbuyingBuyCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockbuyingBuyCall) DoAndReturn(f func(context.Context, int64, string) error) *MockbuyingBuyCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
