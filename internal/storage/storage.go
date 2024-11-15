package storage

import (
	"errors"
	"log/slog"

	"github.com/gintokos/tasksrestapi/internal/domain/models"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	GetAllTasks(logger *slog.Logger) ([]models.Task, error)
	CreateTask(task models.Task, logger *slog.Logger) (models.Task, error)
	UpdateTask(task models.Task, logger *slog.Logger) (models.Task, error)
	DeleteTask(id int64, logger *slog.Logger) error
}
