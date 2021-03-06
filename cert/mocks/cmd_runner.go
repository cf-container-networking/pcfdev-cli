// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/cert (interfaces: CmdRunner)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of CmdRunner interface
type MockCmdRunner struct {
	ctrl     *gomock.Controller
	recorder *_MockCmdRunnerRecorder
}

// Recorder for MockCmdRunner (not exported)
type _MockCmdRunnerRecorder struct {
	mock *MockCmdRunner
}

func NewMockCmdRunner(ctrl *gomock.Controller) *MockCmdRunner {
	mock := &MockCmdRunner{ctrl: ctrl}
	mock.recorder = &_MockCmdRunnerRecorder{mock}
	return mock
}

func (_m *MockCmdRunner) EXPECT() *_MockCmdRunnerRecorder {
	return _m.recorder
}

func (_m *MockCmdRunner) Run(_param0 string, _param1 ...string) ([]byte, error) {
	_s := []interface{}{_param0}
	for _, _x := range _param1 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Run", _s...)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdRunnerRecorder) Run(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0}, arg1...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Run", _s...)
}
