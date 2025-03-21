/*
 * @Author: yujiajie
 * @Date: 2025-01-15 16:04:21
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-15 16:04:45
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/etcdclient.go
 * @Description:
 */
package internal

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type EtcdClient interface {
	ActiveConnection() *grpc.ClientConn
	Close() error
	Ctx() context.Context
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error)
	KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error)
	Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error)
	Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan
}
