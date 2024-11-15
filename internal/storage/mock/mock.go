package mocks

import (
	"sync"

	"log/slog"

	"github.com/gintokos/tasksrestapi/internal/domain/models"
	"github.com/gintokos/tasksrestapi/internal/lib/id"
	"github.com/gintokos/tasksrestapi/internal/storage"
)

type MockStorage struct {
	mu         sync.Mutex
	tasks      []models.Task
	GetAllFunc func(logger *slog.Logger) ([]models.Task, error)
	CreateFunc func(task models.Task, logger *slog.Logger) (models.Task, error)
	UpdateFunc func(task models.Task, logger *slog.Logger) (models.Task, error)
	DeleteFunc func(id int64, logger *slog.Logger) error
}

func NewMockStorage(tasks []models.Task) *MockStorage {
	return &MockStorage{
		tasks: tasks,
	}
}

func (m *MockStorage) GetAllTasks(logger *slog.Logger) ([]models.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(logger)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.tasks, nil
}

func (m *MockStorage) CreateTask(task models.Task, logger *slog.Logger) (models.Task, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(task, logger)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	task.ID = id.GenerateRandomID()
	m.tasks = append(m.tasks, task)
	return task, nil
}

func (m *MockStorage) UpdateTask(task models.Task, logger *slog.Logger) (models.Task, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(task, logger)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	exists := false
	for _, t := range m.tasks {
		if t.ID == task.ID {
			exists = true
			task.Description = t.Description
			task.Title = t.Title
			task.DueDate = t.DueDate
			task.OverDue = t.OverDue
		}
	}
	if !exists {
		return models.Task{}, storage.ErrNotFound
	}

	return task, nil
}

func (m *MockStorage) DeleteTask(id int64, logger *slog.Logger) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id, logger)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	exists := false
	buff := make([]models.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		if t.ID == id {
			exists = true
			continue
		}
		buff = append(buff, t)
	}
	if !exists {
		return storage.ErrNotFound
	}
	m.tasks = buff

	return nil
}
