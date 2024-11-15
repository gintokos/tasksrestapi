package sqllite

import (
	"database/sql"
	"log/slog"

	"github.com/gintokos/tasksrestapi/internal/domain/models"
	"github.com/gintokos/tasksrestapi/internal/lib/id"
	"github.com/gintokos/tasksrestapi/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagepath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagepath)
	if err != nil {
		return nil, err
	}

	query, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER PRIMARY KEY,
		title TEXT,
		description TEXT,
		due_date DATETIME,
		overdue BOOLEAN
	)
	`)
	if err != nil {
		return nil, err
	}

	_, err = query.Exec()
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (st *Storage) GetAllTasks(logger *slog.Logger) ([]models.Task, error) {
	logger.Info("op: storage.sqllite.GetAllTasks")

	rows, err := st.db.Query("SELECT id, title, description, due_date, overdue FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.OverDue); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (st *Storage) CreateTask(task models.Task, logger *slog.Logger) (models.Task, error) {
	logger.Info("op: storage.sqllite.CreateTask")

	task.ID = id.GenerateRandomID()

	stmt, err := st.db.Prepare("INSERT INTO tasks (id, title, description, due_date, overdue) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return models.Task{}, err
	}
	defer stmt.Close()

	var dueDate interface{}
	if task.DueDate != nil {
		dueDate = *task.DueDate
	} else {
		dueDate = nil
	}

	_, err = stmt.Exec(task.ID, task.Title, task.Description, dueDate, task.OverDue)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (st *Storage) UpdateTask(task models.Task, logger *slog.Logger) (models.Task, error) {
	logger.Info("op: storage.sqllite.UpdateTask")

	exists, err := st.isExistsByID(task.ID)
	if err != nil {
		return models.Task{}, err
	}
	if !exists {
		return models.Task{}, storage.ErrNotFound
	}

	stmt, err := st.db.Prepare("UPDATE tasks SET title = ?, description = ?, due_date = ?, overdue = ? WHERE id = ?")
	if err != nil {
		return models.Task{}, err
	}
	defer stmt.Close()

	var dueDate interface{}
	if task.DueDate != nil {
		dueDate = *task.DueDate
	} else {
		dueDate = nil
	}

	_, err = stmt.Exec(task.Title, task.Description, dueDate, task.OverDue, task.ID)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (st *Storage) DeleteTask(id int64, logger *slog.Logger) error {
	logger.Info("op: storage.sqllite.DeleteTask")

	exists, err := st.isExistsByID(id)
	if err != nil {
		return err
	}
	if !exists {
		return storage.ErrNotFound
	}

	stmt, err := st.db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return  err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) isExistsByID(id int64) (bool, error) {
	var exists bool
	err := st.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
