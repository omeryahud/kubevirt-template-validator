// Automatically generated by MockGen. DO NOT EDIT!
// Source: network.go

package network

import (
	gomock "github.com/golang/mock/gomock"

	v1 "kubevirt.io/client-go/api/v1"
	api "kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

// Mock of NetworkInterface interface
type MockNetworkInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockNetworkInterfaceRecorder
}

// Recorder for MockNetworkInterface (not exported)
type _MockNetworkInterfaceRecorder struct {
	mock *MockNetworkInterface
}

func NewMockNetworkInterface(ctrl *gomock.Controller) *MockNetworkInterface {
	mock := &MockNetworkInterface{ctrl: ctrl}
	mock.recorder = &_MockNetworkInterfaceRecorder{mock}
	return mock
}

func (_m *MockNetworkInterface) EXPECT() *_MockNetworkInterfaceRecorder {
	return _m.recorder
}

func (_m *MockNetworkInterface) Plug(vmi *v1.VirtualMachineInstance, iface *v1.Interface, network *v1.Network, domain *api.Domain, podInterfaceName string) error {
	ret := _m.ctrl.Call(_m, "Plug", vmi, iface, network, domain, podInterfaceName)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockNetworkInterfaceRecorder) Plug(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Plug", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockNetworkInterface) Unplug() {
	_m.ctrl.Call(_m, "Unplug")
}

func (_mr *_MockNetworkInterfaceRecorder) Unplug() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Unplug")
}
