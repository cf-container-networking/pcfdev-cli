// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/vm (interfaces: VBox)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	vbox "github.com/pivotal-cf/pcfdev-cli/vbox"
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

func (_m *MockVBox) ConflictingVMPresent(_param0 string) (bool, error) {
	ret := _m.ctrl.Call(_m, "ConflictingVMPresent", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockVBoxRecorder) ConflictingVMPresent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ConflictingVMPresent", arg0)
}

func (_m *MockVBox) ImportVM(_param0 string) error {
	ret := _m.ctrl.Call(_m, "ImportVM", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVBoxRecorder) ImportVM(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ImportVM", arg0)
}

func (_m *MockVBox) StartVM(_param0 string) (*vbox.VM, error) {
	ret := _m.ctrl.Call(_m, "StartVM", _param0)
	ret0, _ := ret[0].(*vbox.VM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockVBoxRecorder) StartVM(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "StartVM", arg0)
}

func (_m *MockVBox) StopVM(_param0 string) error {
	ret := _m.ctrl.Call(_m, "StopVM", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVBoxRecorder) StopVM(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "StopVM", arg0)
}
