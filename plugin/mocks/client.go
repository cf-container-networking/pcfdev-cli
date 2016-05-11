// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/plugin (interfaces: Client)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	pivnet "github.com/pivotal-cf/pcfdev-cli/pivnet"
)

// Mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *_MockClientRecorder
}

// Recorder for MockClient (not exported)
type _MockClientRecorder struct {
	mock *MockClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &_MockClientRecorder{mock}
	return mock
}

func (_m *MockClient) EXPECT() *_MockClientRecorder {
	return _m.recorder
}

func (_m *MockClient) DownloadOVA(_param0 string) (*pivnet.DownloadReader, error) {
	ret := _m.ctrl.Call(_m, "DownloadOVA", _param0)
	ret0, _ := ret[0].(*pivnet.DownloadReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DownloadOVA(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DownloadOVA", arg0)
}
