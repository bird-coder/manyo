/*
 * @Author: yujiajie
 * @Date: 2024-05-13 17:41:28
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-04 17:27:20
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

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	cfg *config.HttpConfig

	Engine *gin.Engine
	server *http.Server

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

func (s *HttpServer) Start() error {
	var err error
	if len(s.cfg.CertFile) == 0 || len(s.cfg.KeyFile) == 0 {
		err = s.server.ListenAndServe()
	} else {
		err = s.server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("waiting for server(%s) finish...", s.cfg.Addr)
		}
		return err
	}
	return nil
}

func (s *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Info("server(%s) shutdown error: %v", s.cfg.Addr, err)
		return err
	}
	logger.Info("server(%s) shutdown processed success", s.cfg.Addr)
	return nil
}
