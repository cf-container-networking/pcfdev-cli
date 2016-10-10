// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/vm (interfaces: VM)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	vm "github.com/pivotal-cf/pcfdev-cli/vm"
)

// Mock of VM interface
type MockVM struct {
	ctrl     *gomock.Controller
	recorder *_MockVMRecorder
}

// Recorder for MockVM (not exported)
type _MockVMRecorder struct {
	mock *MockVM
}

func NewMockVM(ctrl *gomock.Controller) *MockVM {
	mock := &MockVM{ctrl: ctrl}
	mock.recorder = &_MockVMRecorder{mock}
	return mock
}

func (_m *MockVM) EXPECT() *_MockVMRecorder {
	return _m.recorder
}

func (_m *MockVM) GetDebugLogs() error {
	ret := _m.ctrl.Call(_m, "GetDebugLogs")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) GetDebugLogs() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDebugLogs")
}

func (_m *MockVM) Provision(_param0 *vm.StartOpts) error {
	ret := _m.ctrl.Call(_m, "Provision", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Provision(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Provision", arg0)
}

func (_m *MockVM) Resume() error {
	ret := _m.ctrl.Call(_m, "Resume")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Resume() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Resume")
}

func (_m *MockVM) SSH() error {
	ret := _m.ctrl.Call(_m, "SSH")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) SSH() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SSH")
}

func (_m *MockVM) Start(_param0 *vm.StartOpts) error {
	ret := _m.ctrl.Call(_m, "Start", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Start(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start", arg0)
}

func (_m *MockVM) Status() string {
	ret := _m.ctrl.Call(_m, "Status")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockVMRecorder) Status() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Status")
}

func (_m *MockVM) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}

func (_m *MockVM) Suspend() error {
	ret := _m.ctrl.Call(_m, "Suspend")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Suspend() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Suspend")
}

func (_m *MockVM) Target(_param0 bool) error {
	ret := _m.ctrl.Call(_m, "Target", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Target(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Target", arg0)
}

func (_m *MockVM) Trust(_param0 *vm.StartOpts) error {
	ret := _m.ctrl.Call(_m, "Trust", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) Trust(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Trust", arg0)
}

func (_m *MockVM) VerifyStartOpts(_param0 *vm.StartOpts) error {
	ret := _m.ctrl.Call(_m, "VerifyStartOpts", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockVMRecorder) VerifyStartOpts(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VerifyStartOpts", arg0)
}
