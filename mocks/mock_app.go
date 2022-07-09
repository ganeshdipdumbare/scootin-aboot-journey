// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ganeshdipdumbare/scootin-aboot-journey/app (interfaces: App)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockApp is a mock of App interface.
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp.
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance.
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// BeginTrip mocks base method.
func (m *MockApp) BeginTrip(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTrip", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BeginTrip indicates an expected call of BeginTrip.
func (mr *MockAppMockRecorder) BeginTrip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTrip", reflect.TypeOf((*MockApp)(nil).BeginTrip), arg0, arg1, arg2)
}

// GetNearbyAvailableScooters mocks base method.
func (m *MockApp) GetNearbyAvailableScooters(arg0 context.Context, arg1 domain.GeoLocation, arg2 int) ([]domain.Scooter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNearbyAvailableScooters", arg0, arg1, arg2)
	ret0, _ := ret[0].([]domain.Scooter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNearbyAvailableScooters indicates an expected call of GetNearbyAvailableScooters.
func (mr *MockAppMockRecorder) GetNearbyAvailableScooters(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNearbyAvailableScooters", reflect.TypeOf((*MockApp)(nil).GetNearbyAvailableScooters), arg0, arg1, arg2)
}
