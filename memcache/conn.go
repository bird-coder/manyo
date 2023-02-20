package memcache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	crlf            = []byte("\r\n")
	space           = []byte(" ")
	resultOK        = []byte("OK\r\n")
	resultStored    = []byte("STORED\r\n")
	resultNotStored = []byte("NOT_STORED\r\n")
	resultExists    = []byte("EXISTS\r\n")
	resultNotFound  = []byte("NOT_FOUND\r\n")
	resultDeleted   = []byte("DELETED\r\n")
	resultEnd       = []byte("END\r\n")
	resultOk        = []byte("OK\r\n")
	resultTouched   = []byte("TOUCHED\r\n")

	resultClientErrorPrefix = []byte("CLIENT_ERROR ")
	resultServerErrorPrefix = []byte("SERVER_ERROR ")
	versionPrefix           = []byte("VERSION")
)

func replyToError(line []byte) error {
	switch {
	case bytes.Equal(line, resultStored):
		return nil
	case bytes.Equal(line, resultOK):
		return nil
	case bytes.Equal(line, resultEnd):
		return nil
	case bytes.Equal(line, resultNotStored):
		return ErrNotStored
	case bytes.Equal(line, resultExists):
		return ErrCASConflict
	case bytes.Equal(line, resultNotFound):
		return ErrCacheMiss
	case bytes.HasPrefix(line, resultClientErrorPrefix):
		errMsg := line[len(resultClientErrorPrefix):]
		return protocolError(errMsg)
	case bytes.HasPrefix(line, resultServerErrorPrefix):
		errMsg := line[len(resultServerErrorPrefix):]
		return protocolError(errMsg)
	}
	return protocolError(string(line))
}

func legalKey(key string) bool {
	if len(key) > 250 {
		return false
	}
	for i := 0; i < len(key); i++ {
		if key[i] <= ' ' || key[i] == 0x7f {
			return false
		}
	}
	return true
}

// Item is an item to be got or stored in a memcached server.
type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key string

	// Value is the Item's value.
	Value []byte

	// Delta is the change value of incr/decr command
	Delta int

	// Flags are server-opaque flags whose semantics are entirely
	// up to the app.
	Flags uint32

	// Expiration is the cache expiration time, in seconds: either a relative
	// time from now (up to 1 month), or an absolute Unix epoch time.
	// Zero means the Item has no expiration time.
	Expiration int32

	// Compare and swap ID.
	casid uint64
}

type conn struct {
	mu   sync.Mutex
	err  error
	conn net.Conn

	// Read & Write
	readTimeout  time.Duration
	writeTimeout time.Duration
	rw           *bufio.ReadWriter
}

// DialOption specifies an option for dialing a Memcache server.
type DialOption struct {
	f func(*dialOptions)
}

type dialOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	dialer       *net.Dialer
	protocol     string
	dial         func(network, addr string) (net.Conn, error)
}

// DialReadTimeout specifies the timeout for reading a single command reply.
func DialReadTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.readTimeout = d
	}}
}

// DialWriteTimeout specifies the timeout for writing a single command.
func DialWriteTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.writeTimeout = d
	}}
}

// DialConnectTimeout specifies the timeout for connecting to the Memcache server when
// no DialNetDial option is specified.
func DialConnectTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialer.Timeout = d
	}}
}

// DialKeepAlive specifies the keep-alive period for TCP connections to the Redis server
// when no DialNetDial option is specified.
// If zero, keep-alives are not enabled. If no DialKeepAlive option is specified then
// the default of 5 minutes is used to ensure that half-closed TCP sessions are detected.
func DialKeepAlive(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dialer.KeepAlive = d
	}}
}

// DialNetDial specifies a custom dial function for creating TCP
// connections, otherwise a net.Dialer customized via the other options is used.
// DialNetDial overrides DialConnectTimeout and DialKeepAlive.
func DialNetDial(dial func(network, addr string) (net.Conn, error)) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dial = dial
	}}
}

// Dial connects to the Memcache server at the given network and
// address using the specified options.
func Dial(network, address string, options ...DialOption) (Conn, error) {
	do := dialOptions{
		dialer: &net.Dialer{
			KeepAlive: time.Minute * 5,
		},
	}
	for _, option := range options {
		option.f(&do)
	}
	if do.dial == nil {
		do.dial = do.dialer.Dial
	}

	netConn, err := do.dial(network, address)
	if err != nil {
		return nil, err
	}

	c := &conn{
		conn:         netConn,
		rw:           bufio.NewReadWriter(bufio.NewReader(netConn), bufio.NewWriter(netConn)),
		readTimeout:  do.readTimeout,
		writeTimeout: do.writeTimeout,
	}

	return c, nil
}

func (c *conn) Close() error {
	c.mu.Lock()
	err := c.err
	if c.err == nil {
		c.err = errors.New("memcache: closed")
		err = c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) fatal(err error) error {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		// Close connection to force errors on subsequent calls and to unblock
		// other reader or writer.
		c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *conn) Populate(cmd string, item *Item) error {
	if !legalKey(item.Key) {
		return ErrMalformedKey
	}
	var err error
	c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if cmd == "cas" {
		_, err = fmt.Fprintf(c.rw, "%s %s %d %d %d %d\r\n",
			cmd, item.Key, item.Flags, item.Expiration, len(item.Value), item.casid)
	} else {
		_, err = fmt.Fprintf(c.rw, "%s %s %d %d %d\r\n",
			cmd, item.Key, item.Flags, item.Expiration, len(item.Value))
	}
	if err != nil {
		return c.fatal(err)
	}
	if _, err = c.rw.Write(item.Value); err != nil {
		return c.fatal(err)
	}
	if _, err := c.rw.Write(crlf); err != nil {
		return c.fatal(err)
	}
	if err := c.rw.Flush(); err != nil {
		return c.fatal(err)
	}
	c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return c.fatal(err)
	}
	return replyToError(line)
}

func (c *conn) Get(key string) (result *Item, err error) {
	c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if _, err = fmt.Fprintf(c.rw, "gets %s\r\n", key); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	if err = c.parseGetResponse(func(it *Item) { result = it }); err != nil {
		return
	}
	if result == nil {
		return nil, ErrCacheMiss
	}
	return
}

func (c *conn) GetMulti(keys []string) (map[string]*Item, error) {
	c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if _, err := fmt.Fprintf(c.rw, "gets %s\r\n", strings.Join(keys, " ")); err != nil {
		return nil, c.fatal(err)
	}
	if err := c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	results := make(map[string]*Item, len(keys))
	if err := c.parseGetResponse(func(it *Item) { results[it.Key] = it }); err != nil {
		return nil, err
	}
	return results, nil
}

// parseGetResponse reads a GET response from r and calls cb for each
// read and allocated Item
func (c *conn) parseGetResponse(cb func(*Item)) error {
	c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	for {
		line, err := c.rw.ReadSlice('n')
		if err != nil {
			return c.fatal(err)
		}
		if bytes.Equal(line, resultEnd) {
			return nil
		}
		if bytes.HasPrefix(line, resultClientErrorPrefix) {
			errMsg := line[len(resultClientErrorPrefix):]
			return c.fatal(protocolError(errMsg))
		}
		it := new(Item)
		size, err := scanGetResponseLine(line, it)
		if err != nil {
			return c.fatal(err)
		}
		it.Value = make([]byte, size+2)
		_, err = io.ReadFull(c.rw, it.Value)
		if err != nil {
			it.Value = nil
			return c.fatal(err)
		}
		if !bytes.HasSuffix(it.Value, crlf) {
			it.Value = nil
			return c.fatal(protocolError("corrupt get reply, no except CRLF"))
		}
		it.Value = it.Value[:size]
		cb(it)
	}
}

// scanGetResponseLine populates it and returns the declared size of the item.
// It does not read the bytes of the item.
func scanGetResponseLine(line []byte, it *Item) (size int, err error) {
	pattern := "VALUE %s %d %d %d\r\n"
	dest := []interface{}{&it.Key, &it.Flags, &size, &it.casid}
	if bytes.Count(line, space) == 3 {
		pattern = "VALUE %s %d %d\r\n"
		dest = dest[:3]
	}
	n, err := fmt.Sscanf(string(line), pattern, dest...)
	if err != nil || n != len(dest) {
		return -1, fmt.Errorf("memcache: unexpected line in get response: %q", line)
	}
	return size, nil
}

func (c *conn) writeReadLine(format string, args ...interface{}) ([]byte, error) {
	c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	_, err := fmt.Fprintf(c.rw, format, args...)
	if err != nil {
		return nil, c.fatal(err)
	}
	if err := c.rw.Flush(); err != nil {
		return nil, c.fatal(err)
	}
	c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return nil, c.fatal(err)
	}
	return line, nil
}

func (c *conn) Save(cmd string, item *Item) (uint64, error) {
	if !legalKey(item.Key) {
		return 0, ErrMalformedKey
	}
	var err error
	var line []byte
	switch cmd {
	case "delete":
		line, err = c.writeReadLine("%s %s\r\n", cmd, item.Key)
		break
	case "touch":
		line, err = c.writeReadLine("%s %s %d\r\n", cmd, item.Key, item.Expiration)
		break
	case "incr":
	case "decr":
		line, err = c.writeReadLine("%s %s %d\r\n", cmd, item.Key, item.Delta)
		break
	}
	if err != nil {
		return 0, err
	}
	if err = replyToError(line); err != nil {
		return 0, err
	}
	if cmd == "incr" || cmd == "decr" {
		val, err := strconv.ParseUint(string(line[:len(line)-2]), 10, 64)
		if err != nil {
			return 0, err
		}
		return val, nil
	}
	return 0, nil
}
