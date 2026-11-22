package todos

import (
	"Altyn/projectOne/internal/api"
	"Altyn/projectOne/internal/database"
	"Altyn/projectOne/internal/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

func CreateTODO(cfg *api.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req TaskRequest
		ctx := r.Context()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JsonResponseError(w, http.StatusBadRequest, "Couldn't decode the request")
			return
		}

		taskName := req.TaskName
		if strings.TrimSpace(taskName) == "" {
			taskName = "untitled"
		}

		if strings.TrimSpace(req.TaskDesc) == "" {
			utils.JsonResponseError(w, http.StatusBadRequest, "Task description is required")
			return
		}
		req.TaskDesc = strings.TrimSpace(req.TaskDesc)

		createTask := database.CreateTaskParams{
			TaskName: taskName,
			TaskDesc: req.TaskDesc,
		}

		createTask.Deadline = sql.NullTime{
			Time:  req.Deadline,
			Valid: !req.Deadline.IsZero(),
		}

		task, err := cfg.Db.CreateTask(ctx, createTask)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Couldn't create the task")
			return
		}
		utils.JsonResponse(w, http.StatusCreated, "Task created successfully", task)
	})
}
