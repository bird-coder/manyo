package redis

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/bird-coder/manyo/pool"
)

type Pool struct {
	pool.Pool
}

func NewPool(c *RedisConfig) *Pool {
	if c.DialTimeout <= 0 || c.ReadTimeout <= 0 || c.WriteTimeout <= 0 {
		panic("must config redis timeout")
	}
	if c.SlowLog <= 0 {
		c.SlowLog = time.Duration(250 * time.Millisecond)
	}
	ops := []DialOption{
		DialConnectTimeout(time.Duration(c.DialTimeout)),
		DialReadTimeout(time.Duration(c.ReadTimeout)),
		DialWriteTimeout(time.Duration(c.WriteTimeout)),
		DialPassword(c.Password),
		DialDatabase(c.Db),
	}
	p := pool.NewSlice(c.PoolConfig)
	p.Dial = func() (io.Closer, error) {
		c, err := Dial(c.Protocol, c.Addr, ops...)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return &Pool{Pool: p}
}

func (p *Pool) Get(ctx context.Context) Conn {
	c, err := p.Pool.Get(ctx)
	if err != nil {
		return errorConn{err: err}
	}
	c1, _ := c.(Conn)
	return &activeConn{p: p, c: c1}
}

func (p *Pool) Close() error {
	return p.Pool.Close()
}

type activeConn struct {
	p     *Pool
	c     Conn
	state int
}

var (
	sentinel     []byte
	sentinelOnce sync.Once
)

func initSentinel() {
	p := make([]byte, 64)
	if _, err := rand.Read(p); err == nil {
		sentinel = p
	} else {
		h := sha1.New()
		io.WriteString(h, "Oops, rand failed. Use time instead.")
		io.WriteString(h, strconv.FormatInt(time.Now().UnixNano(), 10))
		sentinel = h.Sum(nil)
	}
}

func (ac *activeConn) Close() error {
	c := ac.c
	if c == nil {
		return nil
	}
	ac.c = nil

	if ac.state&MultiState != 0 {
		c.Send("DISCARD")
		ac.state &^= (MultiState | WatchState)
	} else if ac.state&WatchState != 0 {
		c.Send("UNWATCH")
		ac.state &^= WatchState
	}
	if ac.state&SubscribeState != 0 {
		c.Send("UNSUBSCRIBE")
		c.Send("PUNSUBSCRIBE")
		// To detect the end of the message stream, ask the server to echo
		// a sentinel value and read until we see that value.
		sentinelOnce.Do(initSentinel)
		c.Send("ECHO", sentinel)
		c.Flush()
		for {
			p, err := c.Receive()
			if err != nil {
				break
			}
			if p, ok := p.([]byte); ok && bytes.Equal(p, sentinel) {
				ac.state &^= SubscribeState
				break
			}
		}
	}
	_, err := c.Do("")
	ac.p.Pool.Put(c, ac.state != 0 || c.Err() != nil)
	return err
}

func (ac *activeConn) Err() error {
	return ac.c.Err()
}

func (ac *activeConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	ci := LookupCommandInfo(commandName)
	ac.state = (ac.state | ci.Set) &^ ci.Clear
	return ac.c.Do(commandName, args...)
}

func (ac *activeConn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	cwt, ok := ac.c.(ConnWithTimeout)
	if !ok {
		return nil, errTimeoutNotSupported
	}
	ci := LookupCommandInfo(commandName)
	ac.state = (ac.state | ci.Set) &^ ci.Clear
	return cwt.DoWithTimeout(timeout, commandName, args...)
}

func (ac *activeConn) Send(commandName string, args ...interface{}) error {
	ci := LookupCommandInfo(commandName)
	ac.state = (ac.state | ci.Set) &^ ci.Clear
	return ac.c.Send(commandName, args...)
}

func (ac *activeConn) Flush() error {
	return ac.c.Flush()
}

func (ac *activeConn) Receive() (reply interface{}, err error) {
	return ac.c.Receive()
}

func (ac *activeConn) ReceiveWithTimeout(timeout time.Duration) (reply interface{}, err error) {
	cwt, ok := ac.c.(ConnWithTimeout)
	if !ok {
		return nil, errTimeoutNotSupported
	}
	return cwt.ReceiveWithTimeout(timeout)
}

type errorConn struct{ err error }

func (ec errorConn) Do(string, ...interface{}) (interface{}, error) {
	return nil, ec.err
}
func (ec errorConn) DoWithTimeout(time.Duration, string, ...interface{}) (interface{}, error) {
	return nil, ec.err
}
func (ec errorConn) Send(string, ...interface{}) error {
	return ec.err
}
func (ec errorConn) Err() error {
	return ec.err
}
func (ec errorConn) Close() error {
	return ec.err
}
func (ec errorConn) Flush() error {
	return ec.err
}
func (ec errorConn) Receive() (interface{}, error) {
	return nil, ec.err
}
func (ec errorConn) ReceiveWithTimeout(time.Duration) (interface{}, error) {
	return nil, ec.err
}
