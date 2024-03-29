// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/port/paymentprocessor.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/payment-gateway/internal/domain/model"
)

// MockIPaymentProcessor is a mock of IPaymentProcessor interface.
type MockIPaymentProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockIPaymentProcessorMockRecorder
}

// MockIPaymentProcessorMockRecorder is the mock recorder for MockIPaymentProcessor.
type MockIPaymentProcessorMockRecorder struct {
	mock *MockIPaymentProcessor
}

// NewMockIPaymentProcessor creates a new mock instance.
func NewMockIPaymentProcessor(ctrl *gomock.Controller) *MockIPaymentProcessor {
	mock := &MockIPaymentProcessor{ctrl: ctrl}
	mock.recorder = &MockIPaymentProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPaymentProcessor) EXPECT() *MockIPaymentProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockIPaymentProcessor) Process(session model.Session) (*model.TransactionResult, *model.DomainError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", session)
	ret0, _ := ret[0].(*model.TransactionResult)
	ret1, _ := ret[1].(*model.DomainError)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockIPaymentProcessorMockRecorder) Process(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockIPaymentProcessor)(nil).Process), session)
}
