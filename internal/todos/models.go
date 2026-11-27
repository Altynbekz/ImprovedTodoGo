package todos

import (
	"Altyn/projectOne/internal/database"
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	TaskName string    `json:"task_name"`
	TaskDesc string    `json:"task_desc"`
	Deadline time.Time `json:"deadline"`
}

type TaskResponse struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	TaskName  string     `json:"task_name"`
	TaskDesc  string     `json:"task_desc"`
	Deadline  *time.Time `json:"deadline,omitempty"`
	UserID    uuid.UUID  `json:"user_id"`
}

type UserRequest struct {
	FirstName string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserResponse struct {
	FirstName string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore interface {
	FindUserByEmail(context.Context, string) (bool, error)
	FindUserByUsername(context.Context, string) (bool, error)
	CreateUser(context.Context, database.CreateUserParams) (database.User, error)
}

type UserUpdate interface {
	GetUserById(context.Context, uuid.UUID) (database.User, error)
	FindUserByEmail(context.Context, string) (bool, error)
	FindUserByUsername(context.Context, string) (bool, error)
	UpdateUser(context.Context, database.UpdateUserParams) error
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
