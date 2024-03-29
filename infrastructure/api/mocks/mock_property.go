// Code generated by MockGen. DO NOT EDIT.
// Source: ./property.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	model "lahaus/domain/model"
	properties "lahaus/domain/usecases/properties"
	reflect "reflect"
)

// MockPropertyExecutor is a mock of PropertyExecutor interface
type MockPropertyExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockPropertyExecutorMockRecorder
}

// MockPropertyExecutorMockRecorder is the mock recorder for MockPropertyExecutor
type MockPropertyExecutorMockRecorder struct {
	mock *MockPropertyExecutor
}

// NewMockPropertyExecutor creates a new mock instance
func NewMockPropertyExecutor(ctrl *gomock.Controller) *MockPropertyExecutor {
	mock := &MockPropertyExecutor{ctrl: ctrl}
	mock.recorder = &MockPropertyExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPropertyExecutor) EXPECT() *MockPropertyExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockPropertyExecutor) Execute(property *model.Property) (*model.Property, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", property)
	ret0, _ := ret[0].(*model.Property)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockPropertyExecutorMockRecorder) Execute(property interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockPropertyExecutor)(nil).Execute), property)
}

// MockSearchPropertyExecutor is a mock of SearchPropertyExecutor interface
type MockSearchPropertyExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockSearchPropertyExecutorMockRecorder
}

// MockSearchPropertyExecutorMockRecorder is the mock recorder for MockSearchPropertyExecutor
type MockSearchPropertyExecutorMockRecorder struct {
	mock *MockSearchPropertyExecutor
}

// NewMockSearchPropertyExecutor creates a new mock instance
func NewMockSearchPropertyExecutor(ctrl *gomock.Controller) *MockSearchPropertyExecutor {
	mock := &MockSearchPropertyExecutor{ctrl: ctrl}
	mock.recorder = &MockSearchPropertyExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSearchPropertyExecutor) EXPECT() *MockSearchPropertyExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockSearchPropertyExecutor) Execute(search properties.PropertySearchParams) (*model.PropertiesPaging, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", search)
	ret0, _ := ret[0].(*model.PropertiesPaging)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockSearchPropertyExecutorMockRecorder) Execute(search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSearchPropertyExecutor)(nil).Execute), search)
}
