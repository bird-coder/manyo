/*
 * @Author: yujiajie
 * @Date: 2025-01-09 22:35:51
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:03:32
 * @FilePath: /manyo/pkg/zrpc/internal/balancer/p2c/p2c.go
 * @Description:
 */
package p2c

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/zrpc/internal/codes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const (
	// Name is the name of p2c balancer.
	Name = "p2c_ewma"

	// default value from finagle
	decayTime int64 = 10
	// if a node is not selected for a period of time, it is forcibly selected.
	forcePick int64 = 1
	// initial success count
	initSuccess = 1000
	// success count to trigger throttling
	throttleSuccess = initSuccess / 2
	// penalty value for load calculation
	penalty = int64(math.MaxInt32)
	// number of pick attempts
	pickTimes = 3
	// log interval for statistics
	logInterval int64 = 60
)

var emptyPickResult balancer.PickResult

func init() {
	balancer.Register(newBuilder())
}

type p2cPickerBuilder struct{}

func (b *p2cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	readySCs := info.ReadySCs
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	conns := make([]*subConn, 0, len(readySCs))
	for conn, connInfo := range readySCs {
		conns = append(conns, &subConn{
			addr:    connInfo.Address,
			conn:    conn,
			success: initSuccess,
		})
	}

	return &p2cPicker{
		conns: conns,
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		stamp: 0,
	}
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(p2cPickerBuilder), base.Config{HealthCheck: true})
}

type p2cPicker struct {
	conns []*subConn
	r     *rand.Rand
	stamp int64
	lock  sync.Mutex
}

func (p *p2cPicker) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	var chosen *subConn
	switch len(p.conns) {
	case 0:
		return emptyPickResult, balancer.ErrNoSubConnAvailable
	case 1:
		chosen = p.choose(p.conns[0], nil)
	case 2:
		chosen = p.choose(p.conns[0], p.conns[1])
	default:
		var node1, node2 *subConn
		for i := 0; i < pickTimes; i++ {
			a := p.r.Intn(len(p.conns))
			b := p.r.Intn(len(p.conns) - 1)
			if b >= a {
				b++
			}
			node1 = p.conns[a]
			node2 = p.conns[b]
			if node1.healthy() && node2.healthy() {
				break
			}
		}
		chosen = p.choose(node1, node2)
	}

	atomic.AddInt64(&chosen.inflight, 1)
	atomic.AddInt64(&chosen.requests, 1)

	return balancer.PickResult{
		SubConn: chosen.conn,
		Done:    p.buildDoneFunc(chosen),
	}, nil
}

func (p *p2cPicker) buildDoneFunc(c *subConn) func(info balancer.DoneInfo) {
	start := time.Now().Unix()
	return func(info balancer.DoneInfo) {
		atomic.AddInt64(&c.inflight, -1)
		now := time.Now().Unix()
		last := atomic.SwapInt64(&c.last, now)
		td := now - last
		if td < 0 {
			td = 0
		}

		w := math.Exp(float64(-td) / float64(decayTime))
		lag := now - start
		if lag < 0 {
			lag = 0
		}
		olag := atomic.LoadUint64(&c.lag)
		if olag == 0 {
			w = 0
		}

		atomic.StoreUint64(&c.lag, uint64(float64(olag)*w+float64(lag)*(1-w)))
		success := initSuccess
		if info.Err != nil && !codes.Acceptable(info.Err) {
			success = 0
		}
		osucc := atomic.LoadUint64(&c.success)
		atomic.StoreUint64(&c.success, uint64(float64(osucc)*w+float64(success)*(1-w)))

		stamp := atomic.LoadInt64(&p.stamp)
		if now-stamp >= logInterval {
			if atomic.CompareAndSwapInt64(&p.stamp, stamp, now) {
				p.logStats()
			}
		}
	}
}

func (p *p2cPicker) choose(c1, c2 *subConn) *subConn {
	start := time.Now().Unix()
	if c2 == nil {
		atomic.StoreInt64(&c1.pick, start)
		return c1
	}

	if c1.load() > c2.load() {
		c1, c2 = c2, c1
	}

	pick := atomic.LoadInt64(&c2.pick)
	if start-pick > forcePick && atomic.CompareAndSwapInt64(&c2.pick, pick, start) {
		return c2
	}

	atomic.StoreInt64(&c1.pick, start)
	return c1
}

func (p *p2cPicker) logStats() {
	stats := make([]string, 0, len(p.conns))

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, conn := range p.conns {
		stats = append(stats, fmt.Sprintf("conn: %s, load: %d, reqs: %d",
			conn.addr.Addr, conn.load(), atomic.SwapInt64(&conn.requests, 0)))
	}

	logger.Info("p2c - %s", strings.Join(stats, "; "))
}

type subConn struct {
	// The request latency measured by the weighted moving average algorithm.
	lag uint64

	// The value represents the number of requests that are either pending or just
	// starting at the current node, and it is obtained through atomic addition.
	inflight int64
	success  uint64
	requests int64
	last     int64
	pick     int64
	addr     resolver.Address
	conn     balancer.SubConn
}

func (c *subConn) healthy() bool {
	return atomic.LoadUint64(&c.success) > throttleSuccess
}

func (c *subConn) load() int64 {
	// plus one to avoid multiply zero
	lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
	load := lag * (atomic.LoadInt64(&c.inflight) + 1)
	if load == 0 {
		return penalty
	}
	return load
}
