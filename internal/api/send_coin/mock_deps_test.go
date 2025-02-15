// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source deps.go -package send_coin -typed -destination mock_deps_test.go
//

// Package send_coin is a generated GoMock package.
package send_coin

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockcoinSending is a mock of coinSending interface.
type MockcoinSending struct {
	ctrl     *gomock.Controller
	recorder *MockcoinSendingMockRecorder
}

// MockcoinSendingMockRecorder is the mock recorder for MockcoinSending.
type MockcoinSendingMockRecorder struct {
	mock *MockcoinSending
}

// NewMockcoinSending creates a new mock instance.
func NewMockcoinSending(ctrl *gomock.Controller) *MockcoinSending {
	mock := &MockcoinSending{ctrl: ctrl}
	mock.recorder = &MockcoinSendingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcoinSending) EXPECT() *MockcoinSendingMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockcoinSending) Send(ctx context.Context, employeeID int64, targetUsername string, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, employeeID, targetUsername, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockcoinSendingMockRecorder) Send(ctx, employeeID, targetUsername, amount any) *MockcoinSendingSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockcoinSending)(nil).Send), ctx, employeeID, targetUsername, amount)
	return &MockcoinSendingSendCall{Call: call}
}

// MockcoinSendingSendCall wrap *gomock.Call
type MockcoinSendingSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockcoinSendingSendCall) Return(arg0 error) *MockcoinSendingSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockcoinSendingSendCall) Do(f func(context.Context, int64, string, int64) error) *MockcoinSendingSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockcoinSendingSendCall) DoAndReturn(f func(context.Context, int64, string, int64) error) *MockcoinSendingSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
