// Code generated by MockGen. DO NOT EDIT.
// Source: net/http (interfaces: Handler)
//
// Generated by this command:
//
//	mockgen -destination=internal/balancer/mocks/http_handler.go -package=mock_balancer net/http Handler
//

// Package mock_balancer is a generated GoMock package.
package mock_balancer

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
	isgomock struct{}
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// ServeHTTP mocks base method.
func (m *MockHandler) ServeHTTP(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeHTTP", arg0, arg1)
}

// ServeHTTP indicates an expected call of ServeHTTP.
func (mr *MockHandlerMockRecorder) ServeHTTP(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeHTTP", reflect.TypeOf((*MockHandler)(nil).ServeHTTP), arg0, arg1)
}
