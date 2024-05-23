package connpool

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type Slice struct {
	// Dial is an application supplied function for creating and configuring a
	// connection.
	//
	// The connection returned from Dial must not be in a special state
	// (subscribed to pubsub channel, transaction started, ...).
	Dial func() (io.Closer, error)

	chInit uint32

	mu       sync.Mutex
	closed   bool
	active   int
	closeCh  chan struct{}
	freeConn []*poolConn
	reqCh    chan struct{}

	conf *PoolConfig
}

// 创建连接池实例
func NewSlice(c *PoolConfig) *Slice {
	pool := &Slice{
		conf: c,
	}
	pool.startCleanerLocked(time.Duration(pool.conf.IdleTimeout))
	return pool
}

// 重载连接池配置
func (p *Slice) Reload(c *PoolConfig) error {
	p.mu.Lock()
	p.setActive(c.Active)
	p.setIdle(c.Idle)
	p.conf = c
	p.mu.Unlock()
	return nil
}

// 获取当前连接数
func (p *Slice) ActiveCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active
}

// 获取空闲连接数
func (p *Slice) IdleCount() int {
	p.mu.Lock()
	p.mu.Unlock()
	return len(p.freeConn)
}

// 初始化
func (p *Slice) lazyInit() {
	if atomic.LoadUint32(&p.chInit) == 1 {
		return
	}
	p.mu.Lock()
	if p.chInit == 0 {
		p.reqCh = make(chan struct{}, p.conf.Active)
		if p.closed {
			close(p.reqCh)
		} else {
			for i := 0; i < p.conf.Active; i++ {
				p.reqCh <- struct{}{}
			}
		}
		atomic.StoreUint32(&p.chInit, 1)
	}
	p.mu.Unlock()
}

// 获取连接
func (p *Slice) Get(ctx context.Context) (io.Closer, error) {
	if p.conf.Wait && p.conf.Active > 0 {
		p.lazyInit()
		if ctx == nil {
			<-p.reqCh
		} else {
			select {
			case <-p.reqCh:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	p.mu.Lock()

	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	idleTimeout := time.Duration(p.conf.IdleTimeout)

	connNum := len(p.freeConn)
	for connNum > 0 {
		pc := p.popItemLocked()
		p.mu.Unlock()
		if !pc.expired(idleTimeout) {
			return pc.c, nil
		}
		pc.close()
		p.mu.Lock()
		p.active--
		connNum = len(p.freeConn)
	}

	if p.conf.Active > 0 && p.active >= p.conf.Active {
		if !p.conf.Wait {
			p.mu.Unlock()
			return nil, ErrPoolExhausted
		}
	}

	p.active++
	p.mu.Unlock()
	c, err := p.Dial()
	if err != nil {
		p.mu.Lock()
		p.active--
		if p.reqCh != nil && !p.closed {
			p.reqCh <- struct{}{}
		}
		p.mu.Unlock()
		return nil, err
	}

	return c, nil
}

// 连接放回池
func (p *Slice) Put(c io.Closer, forceClose bool) error {
	p.mu.Lock()
	pc := &poolConn{
		c:         c,
		createdAt: nowFunc(),
	}
	if !p.closed && !forceClose {
		p.freeConn = append(p.freeConn, pc)
		connNum := len(p.freeConn)
		if connNum > p.maxIdleItemsLocked() {
			pc = p.popItemLocked()
		} else {
			pc = nil
		}
	}

	if pc != nil {
		p.mu.Unlock()
		pc.close()
		p.mu.Lock()
		p.active--
	}

	if p.reqCh != nil && !p.closed {
		p.reqCh <- struct{}{}
	}
	p.mu.Unlock()
	return nil
}

// 删除链接
func (p *Slice) popItemLocked() *poolConn {
	connNum := len(p.freeConn)
	pc := p.freeConn[0]
	copy(p.freeConn, p.freeConn[1:])
	p.freeConn = p.freeConn[:connNum-1]
	return pc
}

func (p *Slice) setActive(n int) {
	p.conf.Active = n
	if n < 0 {
		p.conf.Active = 0
	}
	syncIdle := p.conf.Active > 0 && p.maxIdleItemsLocked() > p.conf.Active
	if syncIdle {
		p.setIdle(n)
	}
}

func (p *Slice) setIdle(n int) {
	if n > 0 {
		p.conf.Idle = n
	} else {
		p.conf.Idle = -1
	}
	if p.conf.Active > 0 && p.maxIdleItemsLocked() > p.conf.Active {
		p.conf.Idle = p.conf.Active
	}
	var closing []*poolConn
	idleCount := len(p.freeConn)
	maxIdle := p.maxIdleItemsLocked()
	if idleCount > maxIdle {
		closing = p.freeConn[maxIdle:]
		p.freeConn = p.freeConn[:maxIdle]
	}
	for _, pc := range closing {
		pc.close()
	}
}

// 清理空闲连接的协程
func (p *Slice) startCleanerLocked(d time.Duration) {
	if d <= 0 {
		return
	}

	if d < time.Duration(p.conf.IdleTimeout) && p.closeCh != nil {
		select {
		case p.closeCh <- struct{}{}:
		default:
		}
	}

	if p.closeCh == nil {
		p.closeCh = make(chan struct{}, 1)
		go p.staleCleaner()
	}
}

// 清理空闲连接
func (p *Slice) staleCleaner() {
	d := time.Duration(p.conf.IdleTimeout)
	const minInterval = 100 * time.Millisecond
	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
		case <-p.closeCh:
		}
		p.mu.Lock()
		d = time.Duration(p.conf.IdleTimeout)
		if p.closed || d <= 0 {
			p.mu.Unlock()
			return
		}
		var closing []*poolConn
		for i := 0; i < len(p.freeConn); i++ {
			pc := p.freeConn[i]
			if pc.expired(d) {
				closing = append(closing, pc)
				p.active--
				last := len(p.freeConn) - 1
				p.freeConn[i] = p.freeConn[last]
				p.freeConn[last] = nil
				p.freeConn = p.freeConn[:last]
				i--
			}
		}
		p.mu.Unlock()

		for _, pc := range closing {
			pc.close()
		}

		if d < minInterval {
			d = minInterval
		}
		t.Reset(d)
	}
}

const defaultIdleItems = 2

func (p *Slice) maxIdleItemsLocked() int {
	n := p.conf.Idle
	switch {
	case n == 0:
		return defaultIdleItems
	case n < 0:
		return 0
	default:
		return n
	}
}

func (p *Slice) releaseConn(pc *poolConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.freeConn == nil {
		p.freeConn = []*poolConn{}
	}
	p.freeConn = append(p.freeConn, pc)
}

func (p *Slice) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.active -= len(p.freeConn)
	connList := p.freeConn
	p.freeConn = nil
	if p.closeCh != nil {
		close(p.closeCh)
	}
	p.mu.Unlock()
	for _, pc := range connList {
		pc.close()
	}
	return nil
}
