package todos

import (
	"Altyn/projectOne/internal/auth"
	"Altyn/projectOne/internal/database"
	"Altyn/projectOne/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type UserRequest struct {
	FirstName string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserStore interface {
	FindUserByEmail(context.Context, string) (bool, error)
	FindUserByUsername(context.Context, string) (bool, error)
	CreateUser(context.Context, database.CreateUserParams) (database.User, error)
}

func CreateUser(db UserStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		userRequest := &UserRequest{}
		if err := json.NewDecoder(r.Body).Decode(userRequest); err != nil {
			utils.JsonResponseError(w, http.StatusBadRequest, "Could not decode the request")
			return
		}

		userRequest.Username = strings.TrimSpace(userRequest.Username)
		userRequest.Email = strings.TrimSpace(userRequest.Email)
		if userRequest.Username == "" {
			utils.JsonResponseError(w, http.StatusBadRequest, "Username is required")
			return
		}
		if userRequest.Email == "" {
			utils.JsonResponseError(w, http.StatusBadRequest, "Email is required")
			return
		}

		ok, err := db.FindUserByEmail(ctx, userRequest.Email)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Could not check for existing email")
			return
		}

		if ok {
			utils.JsonResponseError(w, http.StatusConflict, "User already exists")
			return
		}

		ok, err = db.FindUserByUsername(ctx, userRequest.Username)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Could not check for existing username")
			return
		}
		if ok {
			utils.JsonResponseError(w, http.StatusConflict, "The username "+userRequest.Username+" already exists")
			return
		}

		if len(strings.TrimSpace(userRequest.Password)) < 8 {
			utils.JsonResponseError(w, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		password, err := auth.HashPassword(userRequest.Password)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Could not hash the password")
			return
		}

		createUser := database.CreateUserParams{
			Useraname: userRequest.Username,
			FirstName: userRequest.FirstName,
			LastName:  userRequest.Lastname,
			Email:     userRequest.Email,
			Password:  password,
		}
		saveUser, err := db.CreateUser(ctx, createUser)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Could not create the user")
			return
		}

		res := database.User{
			ID:        saveUser.ID,
			Useraname: saveUser.Useraname,
			FirstName: saveUser.FirstName,
			LastName:  saveUser.LastName,
			Email:     saveUser.Email,
			CreatedAt: saveUser.CreatedAt,
			UpdatedAt: saveUser.UpdatedAt,
		}

		utils.JsonResponse(w, http.StatusCreated, "Successfully created", res)
	})
}
