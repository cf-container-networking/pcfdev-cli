// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/vbox (interfaces: Address)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Address interface
type MockAddress struct {
	ctrl     *gomock.Controller
	recorder *_MockAddressRecorder
}

// Recorder for MockAddress (not exported)
type _MockAddressRecorder struct {
	mock *MockAddress
}

func NewMockAddress(ctrl *gomock.Controller) *MockAddress {
	mock := &MockAddress{ctrl: ctrl}
	mock.recorder = &_MockAddressRecorder{mock}
	return mock
}

func (_m *MockAddress) EXPECT() *_MockAddressRecorder {
	return _m.recorder
}

func (_m *MockAddress) DomainForIP(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "DomainForIP", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAddressRecorder) DomainForIP(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DomainForIP", arg0)
}

func (_m *MockAddress) SubnetForIP(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "SubnetForIP", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAddressRecorder) SubnetForIP(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SubnetForIP", arg0)
}