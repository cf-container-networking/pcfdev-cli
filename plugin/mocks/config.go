// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/plugin (interfaces: Config)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Config interface
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigRecorder
}

// Recorder for MockConfig (not exported)
type _MockConfigRecorder struct {
	mock *MockConfig
}

func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &_MockConfigRecorder{mock}
	return mock
}

func (_m *MockConfig) EXPECT() *_MockConfigRecorder {
	return _m.recorder
}

func (_m *MockConfig) GetVMName() string {
	ret := _m.ctrl.Call(_m, "GetVMName")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockConfigRecorder) GetVMName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVMName")
}
