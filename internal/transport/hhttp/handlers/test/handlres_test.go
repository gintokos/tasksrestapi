package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gintokos/tasksrestapi/internal/domain/models"
	mocks "github.com/gintokos/tasksrestapi/internal/storage/mock"
	"github.com/gintokos/tasksrestapi/internal/transport/hhttp/handlers"
)

func TestGetTask(t *testing.T) {
	logger := slog.Default()
	mockStorage := mocks.NewMockStorage([]models.Task{
		{ID: 1, Title: "Task 1", Description: "Test Task 1"},
	})

	handler := handlers.GetTask(mockStorage, logger)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var tasks []models.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "Task 1" {
		t.Errorf("expected task title 'Task 1', got '%s'", tasks[0].Title)
	}
}

func TestPostTask(t *testing.T) {
	logger := slog.Default()

	mockStorage := mocks.NewMockStorage([]models.Task{})

	handler := handlers.PostTask(mockStorage, logger)

	newTask := models.Task{
		ID:          1,
		Title:       "New Task",
		Description: "Description for new task",
	}

	bodyValid, _ := json.Marshal(newTask)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(bodyValid))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var createdTask models.Task
	if err := json.NewDecoder(rr.Body).Decode(&createdTask); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	jsonInvalid := []byte(`{
		Title:       "New Task",
		Description: "Description for new task",
		DueDate: "fake date"
		}`)

	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(jsonInvalid))
	rr = httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestPutTask(t *testing.T) {
	logger := slog.Default()

	mockStorage := mocks.NewMockStorage([]models.Task{
		{ID: 1, Title: "Old Task", Description: "Old description"},
	})

	handler := handlers.PutTask(mockStorage, logger)

	updatedTask := models.Task{
		Title:       "Updated Task",
		Description: "Updated description",
	}
	body, _ := json.Marshal(updatedTask)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var modifiedTask models.Task
	if err := json.NewDecoder(rr.Body).Decode(&modifiedTask); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	jsonInvalid := []byte(`{
		Title:       "New Task",
		Description: "Description for new task",
		DueDate: "fake date"
		}`)

	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(jsonInvalid))
	rr = httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	logger := slog.Default()

	mockStorage := mocks.NewMockStorage([]models.Task{
		{ID: 1, Title: "Task to Delete"},
	})

	handler := handlers.DeleteTask(mockStorage, logger)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	tasks, _ := mockStorage.GetAllTasks(logger)
	fmt.Println(tasks)
	if len(tasks) != 0 {
		t.Errorf("expected no tasks left, but some are present")
	}

	req = httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rr = httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	logger := slog.Default()

	mockStorage := mocks.NewMockStorage(nil)

	handler := handlers.GetTask(mockStorage, logger)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
