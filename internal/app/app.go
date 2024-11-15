package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/gintokos/tasksrestapi/internal/app/checker"
	hhttpserver "github.com/gintokos/tasksrestapi/internal/app/hhttp-server"
	"github.com/gintokos/tasksrestapi/internal/config.go"
	"github.com/gintokos/tasksrestapi/internal/lib/logger/sl"
	"github.com/gintokos/tasksrestapi/internal/storage"
)

type App struct {
	checker     checker.Checker
	hhttpserver hhttpserver.HttpServer
	logger      *slog.Logger
}

func NewApp(storage storage.Storage, logger *slog.Logger, cfg config.Config) App {
	return App{
		checker:     checker.NewChecker(logger, storage, cfg.Checker),
		hhttpserver: hhttpserver.NewHttpServer(logger, storage, cfg.Server),
		logger:      logger,
	}
}

func (a *App) MustStart() {
	a.checker.StartCheking()

	err := a.hhttpserver.RunServer()
	if err != nil {
		a.logger.Error("error on starting hhttpserver", sl.Err(err))
		os.Exit(1)
	}
}

func (a *App) GraceFullShutdown(ctx context.Context) error {
	err := a.checker.GraceFullShutdown()
	if err != nil {
		return err
	}
	return a.hhttpserver.GraceFullShutdown(ctx)
}
