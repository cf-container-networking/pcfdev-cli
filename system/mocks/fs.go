// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/system (interfaces: FS)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of FS interface
type MockFS struct {
	ctrl     *gomock.Controller
	recorder *_MockFSRecorder
}

// Recorder for MockFS (not exported)
type _MockFSRecorder struct {
	mock *MockFS
}

func NewMockFS(ctrl *gomock.Controller) *MockFS {
	mock := &MockFS{ctrl: ctrl}
	mock.recorder = &_MockFSRecorder{mock}
	return mock
}

func (_m *MockFS) EXPECT() *_MockFSRecorder {
	return _m.recorder
}

func (_m *MockFS) Read(_param0 string) ([]byte, error) {
	ret := _m.ctrl.Call(_m, "Read", _param0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFSRecorder) Read(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Read", arg0)
}
