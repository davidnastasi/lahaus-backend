// Code generated by MockGen. DO NOT EDIT.
// Source: ./sign_in.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	model "lahaus/domain/model"
	users "lahaus/domain/usecases/users"
	reflect "reflect"
)

// MockStorageManager is a mock of StorageManager interface
type MockStorageManager struct {
	ctrl     *gomock.Controller
	recorder *MockStorageManagerMockRecorder
}

// MockStorageManagerMockRecorder is the mock recorder for MockStorageManager
type MockStorageManagerMockRecorder struct {
	mock *MockStorageManager
}

// NewMockStorageManager creates a new mock instance
func NewMockStorageManager(ctrl *gomock.Controller) *MockStorageManager {
	mock := &MockStorageManager{ctrl: ctrl}
	mock.recorder = &MockStorageManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorageManager) EXPECT() *MockStorageManagerMockRecorder {
	return m.recorder
}

// SaveUser mocks base method
func (m *MockStorageManager) SaveUser(user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUser indicates an expected call of SaveUser
func (mr *MockStorageManagerMockRecorder) SaveUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockStorageManager)(nil).SaveUser), user)
}

// GetUser mocks base method
func (m *MockStorageManager) GetUser(emil string) (*model.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", emil)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUser indicates an expected call of GetUser
func (mr *MockStorageManagerMockRecorder) GetUser(emil interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStorageManager)(nil).GetUser), emil)
}

// GetProperty mocks base method
func (m *MockStorageManager) GetProperty(id int64) (*model.Property, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProperty", id)
	ret0, _ := ret[0].(*model.Property)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProperty indicates an expected call of GetProperty
func (mr *MockStorageManagerMockRecorder) GetProperty(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProperty", reflect.TypeOf((*MockStorageManager)(nil).GetProperty), id)
}

// AddFavourite mocks base method
func (m *MockStorageManager) AddFavourite(email, propertyID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavourite", email, propertyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFavourite indicates an expected call of AddFavourite
func (mr *MockStorageManagerMockRecorder) AddFavourite(email, propertyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavourite", reflect.TypeOf((*MockStorageManager)(nil).AddFavourite), email, propertyID)
}

// ListFavourites mocks base method
func (m *MockStorageManager) ListFavourites(search users.FavouritesSearchParams) (*model.PropertiesPaging, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFavourites", search)
	ret0, _ := ret[0].(*model.PropertiesPaging)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFavourites indicates an expected call of ListFavourites
func (mr *MockStorageManagerMockRecorder) ListFavourites(search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFavourites", reflect.TypeOf((*MockStorageManager)(nil).ListFavourites), search)
}
