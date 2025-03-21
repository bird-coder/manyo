/*
 * @Author: yujiajie
 * @Date: 2025-01-23 16:32:08
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 16:37:24
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/statewatcher.go
 * @Description:
 */
package internal

import (
	"context"
	"sync"

	"google.golang.org/grpc/connectivity"
)

type (
	etcdConn interface {
		GetState() connectivity.State
		WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
	}

	stateWatcher struct {
		disconnected bool
		currentState connectivity.State
		listeners    []func()
		lock         sync.Mutex
	}
)

func newStateWatcher() *stateWatcher {
	return new(stateWatcher)
}

func (sw *stateWatcher) addListener(l func()) {
	sw.lock.Lock()
	sw.listeners = append(sw.listeners, l)
	sw.lock.Unlock()
}

func (sw *stateWatcher) notifyListeners() {
	sw.lock.Lock()
	defer sw.lock.Unlock()

	for _, l := range sw.listeners {
		l()
	}
}

func (sw *stateWatcher) updateState(conn etcdConn) {
	sw.currentState = conn.GetState()
	switch sw.currentState {
	case connectivity.TransientFailure, connectivity.Shutdown:
		sw.disconnected = true
	case connectivity.Ready:
		if sw.disconnected {
			sw.disconnected = false
			sw.notifyListeners()
		}
	}
}

func (sw *stateWatcher) watch(conn etcdConn) {
	sw.currentState = conn.GetState()
	for {
		if conn.WaitForStateChange(context.Background(), sw.currentState) {
			sw.updateState(conn)
		}
	}
}
