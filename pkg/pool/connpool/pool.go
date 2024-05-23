package connpool

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	// ErrPoolExhausted connections are exhausted.
	ErrPoolExhausted = errors.New("pool exhausted")
	// ErrPoolClosed connection pool is closed.
	ErrPoolClosed = errors.New("pool closed")

	// nowFunc returns the current time; it's overridden in tests.
	nowFunc = time.Now
)

type PoolConfig struct {
	// Active number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	Active int
	// Idle number of idle connections in the pool.
	Idle int
	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration
	// If Wait is true and the pool is at the Active limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool
}

type poolConn struct {
	createdAt time.Time
	c         io.Closer
}

func (pc *poolConn) expired(timeout time.Duration) bool {
	if timeout <= 0 {
		return false
	}
	return pc.createdAt.Add(timeout).Before(nowFunc())
}

func (pc *poolConn) close() error {
	return pc.c.Close()
}

// Pool interface.
type Pool interface {
	Get(ctx context.Context) (io.Closer, error)
	Put(c io.Closer, forceClose bool) error
	Close() error
}
