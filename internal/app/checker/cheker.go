package checker

import (
	"log/slog"
	"time"

	"github.com/gintokos/tasksrestapi/internal/config.go"
	"github.com/gintokos/tasksrestapi/internal/lib/logger/sl"
	"github.com/gintokos/tasksrestapi/internal/storage"
)

type Checker struct {
	storage storage.Storage
	logger  *slog.Logger
	delay   int64
}

func NewChecker(logger *slog.Logger, storage storage.Storage, cfg config.CheckerConfig) Checker {
	return Checker{
		storage: storage,
		logger:  logger,
		delay:   cfg.Delay,
	}
}

func (ch *Checker) StartCheking() {
	go func() {
		ticker := time.NewTicker(time.Duration(ch.delay) * time.Second)

		for {
			<-ticker.C
			ch.checkStorage()
		}
	}()
}

func (ch *Checker) GraceFullShutdown() error {
	return ch.checkStorage()
}

func (ch *Checker) checkStorage() error {
	tasks, err := ch.storage.GetAllTasks(ch.logger)
	if err != nil {
		ch.logger.Error("error on getting all tasks", sl.Err(err))
		return err
	}

	for _, task := range tasks {
		if task.DueDate == nil {
			continue
		}
		if task.DueDate.Before(time.Now()) {
			task.OverDue = true
			ch.storage.UpdateTask(task, ch.logger)
		}
	}
	return nil
}
