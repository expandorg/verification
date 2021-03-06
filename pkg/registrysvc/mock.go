// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/registrysvc/request.go

// Package mock_registrysvc is a generated GoMock package.
package registrysvc

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRegistrySVC is a mock of RegistrySVC interface
type MockRegistrySVC struct {
	ctrl     *gomock.Controller
	recorder *MockRegistrySVCMockRecorder
}

// MockRegistrySVCMockRecorder is the mock recorder for MockRegistrySVC
type MockRegistrySVCMockRecorder struct {
	mock *MockRegistrySVC
}

// NewMockRegistrySVC creates a new mock instance
func NewMockRegistrySVC(ctrl *gomock.Controller) *MockRegistrySVC {
	mock := &MockRegistrySVC{ctrl: ctrl}
	mock.recorder = &MockRegistrySVCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRegistrySVC) EXPECT() *MockRegistrySVCMockRecorder {
	return m.recorder
}

// GetRegistration mocks base method
func (m *MockRegistrySVC) GetRegistration(jobID uint64) (*Registration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegistration", jobID)
	ret0, _ := ret[0].(*Registration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegistration indicates an expected call of GetRegistration
func (mr *MockRegistrySVCMockRecorder) GetRegistration(jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegistration", reflect.TypeOf((*MockRegistrySVC)(nil).GetRegistration), jobID)
}
