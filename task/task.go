package task

import "time"

type Task struct {
	ID          int64
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}
