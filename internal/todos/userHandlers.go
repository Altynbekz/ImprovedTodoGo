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
	"strings"
	"time"
)

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

func LoginUserHandler(cfg *api.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		loginUserRequest := &LoginUser{}

		if err := json.NewDecoder(r.Body).Decode(loginUserRequest); err != nil {
			utils.JsonResponseError(w, http.StatusBadRequest, "Couldn't decode the request")
			return
		}

		loginUserRequest.Email = strings.TrimSpace(loginUserRequest.Email)
		if !strings.Contains(loginUserRequest.Email, "@") || strings.Count(loginUserRequest.Email, "@") > 1 {
			utils.JsonResponseError(w, http.StatusBadRequest, "Invalid email")
			return
		}

		if len(strings.TrimSpace(loginUserRequest.Password)) < 8 {
			utils.JsonResponseError(w, http.StatusBadRequest, "The password is invalid")
			return
		}

		user, err := cfg.Db.GetUserByEmail(ctx, loginUserRequest.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.JsonResponseError(w, http.StatusUnauthorized, "Wrong password or email")
				return
			}
			utils.JsonResponseError(w, http.StatusInternalServerError, "Couldn't execute the function")
			return
		}

		ok, err := auth.CheckPassword(loginUserRequest.Password, user.Password)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Couldn't execute the function")
			return
		}

		if !ok {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Wrong password or email")
			return
		}

		token, err := auth.MakeJwt(user.ID, cfg.SecretKey, time.Hour)
		if err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Couldn't create a jwt token")
			return
		}

		res := UserResponse{
			FirstName: user.FirstName,
			Lastname:  user.LastName,
			Username:  user.Useraname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     token,
		}

		utils.JsonResponse(w, http.StatusOK, "User logged in successfully", res)
	})
}

func DeleteUser(cfg *api.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check and validate the user
		username := strings.TrimSpace(r.URL.Query().Get("username"))
		token, err := auth.GetBearer(r.Header)
		if err != nil {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Could not validate the user")
			return
		}

		if username == "" {
			utils.JsonResponseError(w, http.StatusBadRequest, "Enter a username Please")
			return
		}
		//validate the token
		id, err := auth.ValidateJwt(token, cfg.SecretKey)
		if err != nil {
			utils.JsonResponseError(w, http.StatusUnauthorized, "Could not validate the user")
			return
		}

		// get the username by the id
		ctx := r.Context()
		dbUsername, err := cfg.Db.GetUserById(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.JsonResponseError(w, http.StatusUnauthorized, "Could not validate the user")
				return
			}
			utils.JsonResponseError(w, http.StatusInternalServerError, "Could not fetch user")
			return
		}

		if dbUsername.Useraname != username {
			utils.JsonResponseError(w, http.StatusForbidden, "Could not validate the user")
			return
		}

		if err := cfg.Db.DeleteUserById(ctx, dbUsername.ID); err != nil {
			utils.JsonResponseError(w, http.StatusInternalServerError, "Could not delete user")
			return
		}

		utils.RespondWithoutBody(w, http.StatusNoContent)

	})
}
