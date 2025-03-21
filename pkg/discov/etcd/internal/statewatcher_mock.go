/*
 * @Author: yujiajie
 * @Date: 2025-01-23 16:52:02
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 16:59:00
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/statewatcher_mock.go
 * @Description:
 */
package internal

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/connectivity"
)

type MocketcdConn struct {
	ctrl     *gomock.Controller
	recorder *MocketcdConnMockRecorder
}

type MocketcdConnMockRecorder struct {
	mock *MocketcdConn
}

func NewMocketcdConn(ctrl *gomock.Controller) *MocketcdConn {
	mock := &MocketcdConn{ctrl: ctrl}
	mock.recorder = &MocketcdConnMockRecorder{mock: mock}
	return mock
}

func (m *MocketcdConn) EXPECT() *MocketcdConnMockRecorder {
	return m.recorder
}

func (m *MocketcdConn) GetState() connectivity.State {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState")
	ret0, _ := ret[0].(connectivity.State)
	return ret0
}

func (mr *MocketcdConnMockRecorder) GetState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MocketcdConn)(nil).GetState))
}

func (m *MocketcdConn) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForStateChange", ctx, sourceState)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *MocketcdConnMockRecorder) WaitForStateChange(ctx, sourceState any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForStateChange", reflect.TypeOf((*MocketcdConn)(nil).WaitForStateChange), ctx, sourceState)
}
