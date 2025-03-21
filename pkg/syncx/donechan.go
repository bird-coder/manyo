/*
 * @Author: yujiajie
 * @Date: 2025-01-15 15:55:26
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-15 15:57:22
 * @FilePath: /Go-Base/pkg/syncx/donechan.go
 * @Description:
 */
package syncx

import "sync"

type DoneChan struct {
	done chan struct{}
	once sync.Once
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		done: make(chan struct{}),
	}
}

func (dc *DoneChan) Close() {
	dc.once.Do(func() {
		close(dc.done)
	})
}

func (dc *DoneChan) Done() chan struct{} {
	return dc.done
}
