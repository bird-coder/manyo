/*
 * @Author: yujiajie
 * @Date: 2024-05-31 14:42:53
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 16:13:23
 * @FilePath: /manyo/pkg/connection/connection.go
 * @Description:
 */
package connection

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/util/protocol"
)

type Conn struct {
	net.Conn

	cfg *config.ConnConfig

	mu        sync.RWMutex
	inChan    chan []byte
	closeChan chan struct{}
	once      sync.Once
}

func NewConn(conn net.Conn, cfg *config.ConnConfig) *Conn {
	return &Conn{
		Conn:      conn,
		cfg:       cfg,
		closeChan: make(chan struct{}),
	}
}

func (c *Conn) Start() {
	go c.handleMsgLoop()
	c.recvMsgLoop()
}

func (c *Conn) recvMsgLoop() error {
	buf := make([]byte, 0, c.cfg.MaxMsgBytes)
	readBuf := make([]byte, c.cfg.MaxMsgBytes)
	for {
		c.SetReadDeadline(time.Now().Add(time.Duration(c.cfg.ReadTimeout) * time.Second))
		n, err := c.Read(readBuf)
		if err == nil {
			buf = append(buf, readBuf[:n]...)
			buf, err = c.readMsg(buf)
		}
		if err != nil {
			return err
		}
	}
}

func (c *Conn) readMsg(buf []byte) ([]byte, error) {
	var err error
	for len(buf) >= protocol.HeaderLen {

	}
	return buf, err
}

func (c *Conn) handleMsgLoop() {
	for {
		select {
		case msg := <-c.inChan:
			fmt.Println(msg)
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) Send(msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.SetWriteDeadline(time.Now().Add(time.Duration(c.cfg.WriteTimeout) * time.Second))
	c.Write(msg)
}

func (c *Conn) Close() {
	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		close(c.closeChan)
		c.Close()
		close(c.inChan)
	})
}
