/*
 * @Author: yujiajie
 * @Date: 2024-05-13 17:41:28
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-23 17:58:44
 * @FilePath: /manyo/pkg/server/httpx/http.go
 * @Description:
 */
package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/logger"
)

type HttpServer struct {
	*http.Server

	ctx context.Context
}

func NewHttpServer(ctx context.Context, cfg *config.HttpConfig, engine http.Handler) *HttpServer {
	s := &http.Server{
		Addr:           cfg.Addr,
		ReadTimeout:    time.Duration(cfg.ReadTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(cfg.WriteTimeout * int(time.Second)),
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		Handler:        engine,
	}
	srv := &HttpServer{
		Server: s,
		ctx:    ctx,
	}
	return srv
}

func (s *HttpServer) Start() error {
	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Infof("waiting for server(%s) finish...", s.Addr)
		}
		return err
	}
	return nil
}

func (s *HttpServer) Stop() error {
	if err := s.Shutdown(s.ctx); err != nil {
		logger.Infof("server(%s) shutdown error: %v", s.Addr, err)
		return err
	}
	logger.Infof("server(%s) shutdown processed success", s.Addr)
	return nil
}
