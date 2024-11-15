package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gintokos/tasksrestapi/internal/domain/models"
	"github.com/gintokos/tasksrestapi/internal/domain/server"
	"github.com/gintokos/tasksrestapi/internal/lib/id"
	"github.com/gintokos/tasksrestapi/internal/lib/logger/sl"
	"github.com/gintokos/tasksrestapi/internal/services"
	"github.com/gintokos/tasksrestapi/internal/storage"
)

// to do not found eerr

var internalError = "internal error"
var notFound = "not found"

func GetTask(st storage.Storage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "GET.tasks"
		logger.Info(fmt.Sprintf("op: %s", op))
		defer r.Body.Close()

		tasks, err := services.GetallTasks(st, logger)
		if err != nil {
			logger.Error("error on getting alltasks", sl.Err(err))
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}
		if len(tasks) == 0 {
			WriteNewResponceWithError(w, notFound, http.StatusNotFound, logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			logger.Error("error on encoding task to json", sl.Err(err))
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func PostTask(st storage.Storage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "POST.tasks"
		logger.Info(fmt.Sprintf("op: %s", op))
		defer r.Body.Close()

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			logger.Warn("error on decoding body of request", sl.Err(err))
			WriteNewResponceWithError(w, "invalid credentionals", http.StatusBadRequest, logger)
			return
		}

		taskwithid, err := services.CreateNewTask(task, st, logger)
		if err != nil {
			logger.Error("error on creating task", sl.Err(err))
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}

		responseBuffer := &bytes.Buffer{}
		if err := json.NewEncoder(responseBuffer).Encode(taskwithid); err != nil {
			logger.Error("error on encoding taskwithid to json", sl.Err(err))
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if _, err := w.Write(responseBuffer.Bytes()); err != nil {
			logger.Error("error on writing response to client", sl.Err(err))
		}
	}
}
func PutTask(st storage.Storage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "PUT.tasks.id"
		logger.Info(fmt.Sprintf("op: %s", op))
		defer r.Body.Close()

		idstring := strings.TrimPrefix(r.URL.Path, "/tasks/")
		idint64, ok := id.ValidateID(idstring)
		if !ok {
			logger.Info("putted wrong id")
			WriteNewResponceWithError(w, "invalid id", http.StatusBadRequest, logger)
			return
		}

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			logger.Warn("error on decoding body of request", sl.Err(err))
			WriteNewResponceWithError(w, "invalid credentionals", http.StatusBadRequest, logger)
			return
		}
		task.ID = idint64

		modifiedtask, err := services.UpdateTask(task, st, logger)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				logger.Info("Not found record with this id")
				WriteNewResponceWithError(w, notFound, http.StatusBadRequest, logger)
				return
			}
			logger.Error("error on updating task", sl.Err(err))
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}

		if err := json.NewEncoder(w).Encode(modifiedtask); err != nil {
			logger.Error("error on encoding modifiedtask to json", sl.Err(err))
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}
	}
}

func DeleteTask(st storage.Storage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "DELETE.tasks.ID"
		logger.Info(fmt.Sprintf("op: %s", op))
		defer r.Body.Close()

		idstring := strings.TrimPrefix(r.URL.Path, "/tasks/")
		idint64, ok := id.ValidateID(idstring)
		if !ok {
			logger.Info("putted wrong id")
			WriteNewResponceWithError(w, "invalid id", http.StatusBadRequest, logger)
			return
		}

		err := services.DeleteTask(idint64, st, logger)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				logger.Info("Not found record with this id")
				WriteNewResponceWithError(w, notFound, http.StatusBadRequest, logger)
				return
			}
			logger.Info("Not found record with this id")
			WriteNewResponceWithError(w, internalError, http.StatusInternalServerError, logger)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func WriteNewResponceWithError(w http.ResponseWriter, errstring string, status int, logger *slog.Logger) {
	logger.Info("http-server.handlers.WriteNewResponceWithError")
	w.WriteHeader(status)
	resp := server.ResponceWithError{
		Msg: "error",
		Err: errstring,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("error on encoding errresponce to json", sl.Err(err))
	}
}
