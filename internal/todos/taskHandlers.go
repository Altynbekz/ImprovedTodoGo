package todos

import (
	"Altyn/projectOne/internal/api"
	"Altyn/projectOne/internal/auth"
	"Altyn/projectOne/internal/database"
	"Altyn/projectOne/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func CreateTODO(cfg *api.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//get the token and verify
		token, err := auth.GetBearer(r.Header)
		if err != nil {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Couldn't receive the token")
			return
		}

		//validate the token
		userId, err := auth.ValidateJwt(token, cfg.SecretKey)
		if err != nil {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Could not receive any respond after validating")
			return
		}

		var req TaskRequest
		ctx := r.Context()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JsonResponseError(w, http.StatusBadRequest, "Couldn't decode the request")
			return
		}

		taskName := strings.TrimSpace(req.TaskName)
		if taskName == "" {
			names, err := cfg.Db.GetDefaultName(ctx, userId)
			if err != nil {
				utils.JsonResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
				return
			}

			if len(names) == 0 {
				taskName = "untitled"
			} else {
				num, err := GetTheLastUntitled(names)
				if err != nil {
					utils.JsonResponseError(w, http.StatusInternalServerError, "Could not generate task name")
					return
				}

				taskName = "untitled " + strconv.Itoa(num)
			}

		}

		if strings.TrimSpace(req.TaskDesc) == "" {
			utils.JsonResponseError(w, http.StatusBadRequest, "Task description is required")
			return
		}
		req.TaskDesc = strings.TrimSpace(req.TaskDesc)

		createTask := database.CreateTaskParams{
			TaskName: taskName,
			TaskDesc: req.TaskDesc,
			UserID:   uuid.UUID(userId),
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

		utils.JsonResponse(w, http.StatusCreated, "Task was created successfully", taskToResponse(task))
	})
}

func GetAllTodo(cfg *api.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := strings.TrimSpace(r.URL.Query().Get("username"))
		if username == "" {
			utils.JsonResponseError(w, http.StatusBadRequest, "Username is required")
			return
		}

		//validate the users id
		token, err := auth.GetBearer(r.Header)
		if err != nil {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Could not get the token")
			return
		}
		id, err := auth.ValidateJwt(token, cfg.SecretKey)
		if err != nil {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Could not validate the user")
			return
		}
		ctx := r.Context()

		dbUsername, err := cfg.Db.GetUsernameById(ctx, uuid.UUID(id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.JsonResponseError(w, http.StatusUnauthorized, "The user with that id was not found")
				return
			}
			utils.JsonResponseError(w, http.StatusInternalServerError, "The function malfunctioned")
			return
		}

		if username != dbUsername {
			utils.JsonResponseError(w, http.StatusForbidden, "Invalid username")
			return
		}

		tasks, err := cfg.Db.GetTaskByUserId(ctx, uuid.UUID(id))
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Function malfunctioned")
			return
		}

		respond := make([]TaskResponse, 0, len(tasks))
		for _, task := range tasks {
			respond = append(respond, taskToResponse(task))
		}

		utils.JsonResponse(w, http.StatusOK, "The tasks were successfully retrieved", respond)
	})
}

func GetTodoByName() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
