package services

import (
	"log/slog"

	"github.com/gintokos/tasksrestapi/internal/domain/models"
	"github.com/gintokos/tasksrestapi/internal/storage"
)

func GetallTasks(storage storage.Storage, logger *slog.Logger) ([]models.Task, error) {
	return storage.GetAllTasks(logger)
}

func CreateNewTask(task models.Task, storage storage.Storage, logger *slog.Logger) (models.Task, error) {
	return storage.CreateTask(task, logger)
}

func UpdateTask(task models.Task, storage storage.Storage, logger *slog.Logger) (models.Task, error) {
	return storage.UpdateTask(task, logger)
}

func DeleteTask(id int64, storage storage.Storage, logger *slog.Logger) error {
	return storage.DeleteTask(id, logger)
}
