// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/einride/xsens-go (interfaces: SerialPort)

// Package mockxsens is a generated GoMock package.
package mockxsens

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockSerialPort is a mock of SerialPort interface
type MockSerialPort struct {
	ctrl     *gomock.Controller
	recorder *MockSerialPortMockRecorder
}

// MockSerialPortMockRecorder is the mock recorder for MockSerialPort
type MockSerialPortMockRecorder struct {
	mock *MockSerialPort
}

// NewMockSerialPort creates a new mock instance
func NewMockSerialPort(ctrl *gomock.Controller) *MockSerialPort {
	mock := &MockSerialPort{ctrl: ctrl}
	mock.recorder = &MockSerialPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSerialPort) EXPECT() *MockSerialPortMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockSerialPort) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockSerialPortMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSerialPort)(nil).Close))
}

// Read mocks base method
func (m *MockSerialPort) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockSerialPortMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockSerialPort)(nil).Read), arg0)
}

// SetReadDeadline mocks base method
func (m *MockSerialPort) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline
func (mr *MockSerialPortMockRecorder) SetReadDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockSerialPort)(nil).SetReadDeadline), arg0)
}

// SetWriteDeadline mocks base method
func (m *MockSerialPort) SetWriteDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline
func (mr *MockSerialPortMockRecorder) SetWriteDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*MockSerialPort)(nil).SetWriteDeadline), arg0)
}

// Write mocks base method
func (m *MockSerialPort) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockSerialPortMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockSerialPort)(nil).Write), arg0)
}