package todos

import (
	"Altyn/projectOne/internal/database"
	"strconv"
	"strings"
	"time"
)

func GetTheLastUntitled(names []string) (int, error) {
	maxNum := 0
	for _, name := range names {
		num, err := GetLastNumber(name)
		if err != nil {
			continue
		}
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1, nil
}

func GetLastNumber(name string) (int, error) {
	suffix := strings.TrimSpace(strings.TrimPrefix(name, "untitled"))
	if suffix == "" {
		return 0, nil
	}
	return strconv.Atoi(suffix)
}

func taskToResponse(task database.Task) TaskResponse {
	var deadline *time.Time
	if task.Deadline.Valid {
		deadline = &task.Deadline.Time
	}

	return TaskResponse{
		ID:        task.ID,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
		TaskName:  task.TaskName,
		TaskDesc:  task.TaskDesc,
		Deadline:  deadline,
		UserID:    task.UserID,
	}
}
