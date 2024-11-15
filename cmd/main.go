package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gintokos/tasksrestapi/internal/app"
	"github.com/gintokos/tasksrestapi/internal/config.go"
	"github.com/gintokos/tasksrestapi/internal/lib/logger/sl"
	sqllite "github.com/gintokos/tasksrestapi/internal/storage/sqlLite"
)

func main() {
	cfg := config.MustLoad("config.json")
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	))
	log.Info("Config and logger was inited")

	storage, err := sqllite.NewStorage(cfg.Sql.Storagepath)
	if err != nil {
		log.Error("error on getting storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Storage was inited")

	app := app.NewApp(storage, log, cfg)
	go app.MustStart()
	log.Info("App have started his work")

	canceled := make(chan os.Signal, 1)
	signal.Notify(canceled, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-canceled

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = app.GraceFullShutdown(ctx)
	if err != nil {
		log.Error("error on shutdowning app", sl.Err(err))
	} else {
		log.Info("App stopped his work succesfully")
	}
}
