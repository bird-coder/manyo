package redis

import (
	"context"
	"errors"
	"time"

	"github.com/bird-coder/common/pool"
)

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }

// Conn represents a connection to a Redis server.
type Conn interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value when the connection is not usable.
	Err() error

	// Do sends a command to the server and returns the received reply.
	Do(commandName string, args ...interface{}) (reply interface{}, err error)

	// Send writes the command to the client's output buffer.
	Send(commandName string, args ...interface{}) error

	// Flush flushes the output buffer to the Redis server.
	Flush() error

	// Receive receives a single reply from the Redis server
	Receive() (reply interface{}, err error)
}

// Argument is the interface implemented by an object which wants to control how
// the object is converted to Redis bulk strings.
type Argument interface {
	// RedisArg returns a value to be encoded as a bulk string per the
	// conversions listed in the section 'Executing Commands'.
	// Implementations should typically return a []byte or string.
	RedisArg() interface{}
}

// Scanner is implemented by an object which wants to control its value is
// interpreted when read from Redis.
type Scanner interface {
	// RedisScan assigns a value from a Redis value. The argument src is one of
	// the reply types listed in the section `Executing Commands`.
	//
	// An error should be returned if the value cannot be stored without
	// loss of information.
	RedisScan(src interface{}) error
}

// ConnWithTimeout is an optional interface that allows the caller to override
// a connection's default read timeout. This interface is useful for executing
// the BLPOP, BRPOP, BRPOPLPUSH, XREAD and other commands that block at the
// server.
//
// A connection's default read timeout is set with the DialReadTimeout dial
// option. Applications should rely on the default timeout for commands that do
// not block at the server.
//
// All of the Conn implementations in this package satisfy the ConnWithTimeout
// interface.
//
// Use the DoWithTimeout and ReceiveWithTimeout helper functions to simplify
// use of this interface.
type ConnWithTimeout interface {
	Conn

	// Do sends a command to the server and returns the received reply.
	// The timeout overrides the read timeout set when dialing the
	// connection.
	DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error)

	// Receive receives a single reply from the Redis server. The timeout
	// overrides the read timeout set when dialing the connection.
	ReceiveWithTimeout(timeout time.Duration) (reply interface{}, err error)
}

var errTimeoutNotSupported = errors.New("redis: connection does not support ConnWithTimeout")

// DoWithTimeout executes a Redis command with the specified read timeout. If
// the connection does not satisfy the ConnWithTimeout interface, then an error
// is returned.
func DoWithTimeout(c Conn, timeout time.Duration, cmd string, args ...interface{}) (interface{}, error) {
	cwt, ok := c.(ConnWithTimeout)
	if !ok {
		return nil, errTimeoutNotSupported
	}
	return cwt.DoWithTimeout(timeout, cmd, args...)
}

// ReceiveWithTimeout receives a reply with the specified read timeout. If the
// connection does not satisfy the ConnWithTimeout interface, then an error is
// returned.
func ReceiveWithTimeout(c Conn, timeout time.Duration) (interface{}, error) {
	cwt, ok := c.(ConnWithTimeout)
	if !ok {
		return nil, errTimeoutNotSupported
	}
	return cwt.ReceiveWithTimeout(timeout)
}

type RedisConfig struct {
	*pool.PoolConfig

	Protocol     string
	Addr         string
	Db           int
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	SlowLog      time.Duration
}

type Redis struct {
	pool *Pool
	conf *RedisConfig
}

func NewRedis(c *RedisConfig) *Redis {
	return &Redis{
		pool: NewPool(c),
		conf: c,
	}
}

func (r *Redis) Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get(ctx)
	defer conn.Close()
	reply, err = conn.Do(commandName, args...)
	return
}

func (r *Redis) Close() error {
	return r.pool.Close()
}

func (r *Redis) Conn(ctx context.Context) Conn {
	return r.pool.Get(ctx)
}
