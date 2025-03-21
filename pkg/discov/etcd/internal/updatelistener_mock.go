/*
 * @Author: yujiajie
 * @Date: 2025-01-23 16:44:41
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 16:50:45
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/updatelistener_mock.go
 * @Description:
 */
package internal

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

type MockUpdateListener struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateListenerMockRecorder
}

type MockUpdateListenerMockRecorder struct {
	mock *MockUpdateListener
}

func NewMockUpdateListener(ctrl *gomock.Controller) *MockUpdateListener {
	mock := &MockUpdateListener{ctrl: ctrl}
	mock.recorder = &MockUpdateListenerMockRecorder{mock: mock}
	return mock
}

func (m *MockUpdateListener) EXPECT() *MockUpdateListenerMockRecorder {
	return m.recorder
}

func (m *MockUpdateListener) OnAdd(kv KV) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnAdd", kv)
}

func (mr *MockUpdateListenerMockRecorder) OnAdd(kv any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnAdd", reflect.TypeOf((*MockUpdateListener)(nil).OnAdd), kv)
}

func (m *MockUpdateListener) OnDelete(kv KV) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDelete", kv)
}

func (mr *MockUpdateListenerMockRecorder) OnDelete(kv any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDelete", reflect.TypeOf((*MockUpdateListener)(nil).OnDelete), kv)
}
