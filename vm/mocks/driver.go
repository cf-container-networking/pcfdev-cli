// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/vm (interfaces: Driver)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Driver interface
type MockDriver struct {
	ctrl     *gomock.Controller
	recorder *_MockDriverRecorder
}

// Recorder for MockDriver (not exported)
type _MockDriverRecorder struct {
	mock *MockDriver
}

func NewMockDriver(ctrl *gomock.Controller) *MockDriver {
	mock := &MockDriver{ctrl: ctrl}
	mock.recorder = &_MockDriverRecorder{mock}
	return mock
}

func (_m *MockDriver) EXPECT() *_MockDriverRecorder {
	return _m.recorder
}

func (_m *MockDriver) GetHostForwardPort(_param0 string, _param1 string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetHostForwardPort", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDriverRecorder) GetHostForwardPort(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetHostForwardPort", arg0, arg1)
}

func (_m *MockDriver) GetVMIP(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetVMIP", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDriverRecorder) GetVMIP(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVMIP", arg0)
}

func (_m *MockDriver) VMExists(_param0 string) (bool, error) {
	ret := _m.ctrl.Call(_m, "VMExists", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDriverRecorder) VMExists(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VMExists", arg0)
}

func (_m *MockDriver) VMState(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "VMState", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDriverRecorder) VMState(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VMState", arg0)
}
