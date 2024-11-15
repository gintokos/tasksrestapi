package models

import "time"

type Task struct {
	ID          int64
	Title       string
	Description string
	DueDate     *time.Time
	OverDue     bool
}
