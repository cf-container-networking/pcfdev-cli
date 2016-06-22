// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/plugin (interfaces: VBox)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	config "github.com/pivotal-cf/pcfdev-cli/config"
)

// Mock of VBox interface
type MockVBox struct {
	ctrl     *gomock.Controller
	recorder *_MockVBoxRecorder
}

// Recorder for MockVBox (not exported)
type _MockVBoxRecorder struct {
	mock *MockVBox
}

func NewMockVBox(ctrl *gomock.Controller) *MockVBox {
	mock := &MockVBox{ctrl: ctrl}
	mock.recorder = &_MockVBoxRecorder{mock}
	return mock
}

func (_m *MockVBox) EXPECT() *_MockVBoxRecorder {
	return _m.recorder
}

func (_m *MockVBox) ConflictingVMPresent(_param0 *config.VMConfig) (bool, error) {
	ret := _m.ctrl.Call(_m, "ConflictingVMPresent", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockVBoxRecorder) ConflictingVMPresent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ConflictingVMPresent", arg0)
}

func (_m *MockVBox) DestroyPCFDevVMs() error {
	ret := _m.ctrl.Call(_m, "DestroyPCFDevVMs")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVBoxRecorder) DestroyPCFDevVMs() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DestroyPCFDevVMs")
}

func (_m *MockVBox) GetVMName() (string, error) {
	ret := _m.ctrl.Call(_m, "GetVMName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockVBoxRecorder) GetVMName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVMName")
}
