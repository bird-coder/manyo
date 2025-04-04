/*
 * @Author: yujiajie
 * @Date: 2025-01-05 21:29:18
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-05 21:32:36
 * @FilePath: /Go-Base/pkg/syncx/atomicduration.go
 * @Description:
 */
package syncx

import (
	"sync/atomic"
	"time"
)

type AtomicDuration int64

func NewAtomicDuration() *AtomicDuration {
	return new(AtomicDuration)
}

func ForAtomicDuration(val time.Duration) *AtomicDuration {
	d := NewAtomicDuration()
	d.Set(val)
	return d
}

func (d *AtomicDuration) CompareAndSwap(old, val time.Duration) bool {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(old), int64(val))
}

func (d *AtomicDuration) Load() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

func (d *AtomicDuration) Set(val time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(val))
}
