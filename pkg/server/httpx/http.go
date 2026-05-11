/*
 * @Author: yujiajie
 * @Date: 2024-05-13 17:41:28
 * @LastEditors: yujiajie 1037297660@qq.com
 * @LastEditTime: 2026-05-09 10:52:37
 * @FilePath: /manyo/pkg/server/httpx/http.go
 * @Description:
 */
package httpx

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/logger"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	cfg *config.HttpConfig

	Engine *gin.Engine
	server *http.Server
	ln     net.Listener

	ctx context.Context
}

func NewHttpServer(ctx context.Context, cfg *config.HttpConfig) *HttpServer {
	srv := &HttpServer{
		cfg:    cfg,
		Engine: gin.New(),
		ctx:    ctx,
	}
	srv.init()

	return srv
}

func (s *HttpServer) init() {
	s.server = &http.Server{
		Addr:           s.cfg.Addr,
		ReadTimeout:    time.Duration(s.cfg.ReadTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(s.cfg.WriteTimeout * int(time.Second)),
		MaxHeaderBytes: s.cfg.MaxHeaderBytes,
		Handler:        s.Engine,
	}
}

func (s *HttpServer) Prepare(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	logger.Info("server(%s) is listening", s.cfg.Addr)

	return nil
}

func (s *HttpServer) Start(ctx context.Context) error {
	var err error
	if len(s.cfg.CertFile) == 0 || len(s.cfg.KeyFile) == 0 {
		err = s.server.Serve(s.ln)
	} else {
		err = s.server.ServeTLS(s.ln, s.cfg.CertFile, s.cfg.KeyFile)
	}

	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("server(%s) stopped", s.cfg.Addr)
		return nil
	}
	return err
}

func (s *HttpServer) Stop(ctx context.Context) error {
	if ctx == nil {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Error("server(%s) shutdown error: %v", s.cfg.Addr, err)
		return err
	}
	logger.Info("server(%s) shutdown processed success", s.cfg.Addr)
	return nil
}
