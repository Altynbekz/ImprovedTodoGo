package todos

import (
	"time"
)

type TaskRequest struct {
	TaskName string    `json:"task_name"`
	TaskDesc string    `json:"task_desc"`
	Deadline time.Time `json:"deadline"`
}
