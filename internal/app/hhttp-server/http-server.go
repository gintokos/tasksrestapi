package hhttpserver

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gintokos/tasksrestapi/internal/config.go"
	"github.com/gintokos/tasksrestapi/internal/storage"
	"github.com/gintokos/tasksrestapi/internal/transport/hhttp"
)

type HttpServer struct {
	storage storage.Storage
	logger  *slog.Logger
	server  *http.Server
}

func NewHttpServer(logger *slog.Logger, storage storage.Storage, cfg config.ServerConfig) HttpServer {
	srv := http.Server{
		Addr: "0.0.0.0:8080",
		ErrorLog:          log.New(io.Discard, "", 0),
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
	}
	return HttpServer{
		server:  &srv,
		storage: storage,
		logger:  logger,
	}
}

func (s *HttpServer) RunServer() error {
	router := hhttp.NewRouter(s.storage, s.logger)

	s.server.Handler = router

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HttpServer) GraceFullShutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
