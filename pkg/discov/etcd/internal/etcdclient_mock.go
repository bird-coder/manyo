package internal

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type MockEtcdClient struct {
	ctrl     *gomock.Controller
	recorder *MockEtcdClientMockRecorder
}

type MockEtcdClientMockRecorder struct {
	mock *MockEtcdClient
}

func NewMockEtcdClient(ctrl *gomock.Controller) *MockEtcdClient {
	mock := &MockEtcdClient{ctrl: ctrl}
	mock.recorder = &MockEtcdClientMockRecorder{mock: mock}
	return mock
}

func (m *MockEtcdClient) EXPECT() *MockEtcdClientMockRecorder {
	return m.recorder
}

func (m *MockEtcdClient) ActiveConnection() *grpc.ClientConn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActiveConnection")
	ret0, _ := ret[0].(*grpc.ClientConn)
	return ret0
}

func (mr *MockEtcdClientMockRecorder) ActiveConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActiveConnection", reflect.TypeOf((*MockEtcdClient)(nil).ActiveConnection))
}

func (m *MockEtcdClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockEtcdClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockEtcdClient)(nil).Close))
}

func (m *MockEtcdClient) Ctx() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ctx")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

func (mr *MockEtcdClientMockRecorder) Ctx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ctx", reflect.TypeOf((*MockEtcdClient)(nil).Ctx))
}

func (m *MockEtcdClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*clientv3.GetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockEtcdClientMockRecorder) Get(ctx, key any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockEtcdClient)(nil).Get), varargs...)
}

func (m *MockEtcdClient) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Grant", ctx, ttl)
	ret0, _ := ret[0].(*clientv3.LeaseGrantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockEtcdClientMockRecorder) Grant(ctx, ttl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Grant", reflect.TypeOf((*MockEtcdClient)(nil).Grant), ctx, ttl)
}

func (m *MockEtcdClient) KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepAlive", ctx, id)
	ret0, _ := ret[0].(<-chan *clientv3.LeaseKeepAliveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockEtcdClientMockRecorder) KeepAlive(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepAlive", reflect.TypeOf((*MockEtcdClient)(nil).KeepAlive), ctx, id)
}

func (m *MockEtcdClient) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key, val}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Put", varargs...)
	ret0, _ := ret[0].(*clientv3.PutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockEtcdClientMockRecorder) Put(ctx, key, val any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key, val}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockEtcdClient)(nil).Put), varargs...)
}

func (m *MockEtcdClient) Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Revoke", ctx, id)
	ret0, _ := ret[0].(*clientv3.LeaseRevokeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockEtcdClientMockRecorder) Revoke(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revoke", reflect.TypeOf((*MockEtcdClient)(nil).Revoke), ctx, id)
}

func (m *MockEtcdClient) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Watch", varargs...)
	ret0, _ := ret[0].(clientv3.WatchChan)
	return ret0
}

func (mr *MockEtcdClientMockRecorder) Watch(ctx, key any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockEtcdClient)(nil).Watch), varargs...)
}
