package hhttp

import (
	"log/slog"
	"net/http"

	"github.com/gintokos/tasksrestapi/internal/storage"
	"github.com/gintokos/tasksrestapi/internal/transport/hhttp/handlers"
)

func NewRouter(storage storage.Storage, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handlers.GetTask(storage, logger))
	mux.HandleFunc("POST /tasks", handlers.PostTask(storage, logger))
	mux.HandleFunc("PUT /tasks/{id}", handlers.PutTask(storage, logger))
	mux.HandleFunc("DELETE /tasks/{id}", handlers.DeleteTask(storage, logger))

	return mux
}
