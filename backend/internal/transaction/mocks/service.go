// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/transaction/service.go

// Package mock_transaction is a generated GoMock package.
package mock_transaction

import (
	transaction "affiliates-backoffice-backend/internal/transaction"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceI is a mock of ServiceI interface.
type MockServiceI struct {
	ctrl     *gomock.Controller
	recorder *MockServiceIMockRecorder
}

// MockServiceIMockRecorder is the mock recorder for MockServiceI.
type MockServiceIMockRecorder struct {
	mock *MockServiceI
}

// NewMockServiceI creates a new mock instance.
func NewMockServiceI(ctrl *gomock.Controller) *MockServiceI {
	mock := &MockServiceI{ctrl: ctrl}
	mock.recorder = &MockServiceIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceI) EXPECT() *MockServiceIMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockServiceI) Get(ctx context.Context, affiliateID string) (*[]transaction.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, affiliateID)
	ret0, _ := ret[0].(*[]transaction.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockServiceIMockRecorder) Get(ctx, affiliateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockServiceI)(nil).Get), ctx, affiliateID)
}
