package todos

import (
	"Altyn/projectOne/internal/database"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockUserStore struct {
	emailExists    bool
	usernameExists bool
	createUser     database.User
	createErr      error
}

func (m *mockUserStore) FindUserByEmail(context.Context, string) (bool, error) {
	return m.emailExists, nil
}

func (m *mockUserStore) FindUserByUsername(context.Context, string) (bool, error) {
	return m.usernameExists, nil
}

func (m *mockUserStore) CreateUser(context.Context, database.CreateUserParams) (database.User, error) {
	return m.createUser, m.createErr
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		store      *mockUserStore
		wantStatus int
	}{
		{
			name:       "malformed JSON",
			body:       "{",
			store:      &mockUserStore{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: `{
				"first_name":"Altyn",
				"last_name":"Bek",
				"username":"altyn",
				"email":"altyn@example.com",
				"password":"password123"
			}`,
			store:      &mockUserStore{emailExists: true},
			wantStatus: http.StatusConflict,
		},
		{
			name: "password is too short",
			body: `{
				"first_name":"Altyn",
				"last_name":"Bek",
				"username":"altyn",
				"email":"altyn@example.com",
				"password":"short"
			}`,
			store:      &mockUserStore{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "successful creation",
			body: `{
				"first_name":"Altyn",
				"last_name":"Bek",
				"username":"altyn",
				"email":"altyn@example.com",
				"password":"password123"
			}`,
			store: &mockUserStore{
				createUser: database.User{
					ID:        uuid.New(),
					Useraname: "altyn",
					FirstName: "Altyn",
					LastName:  "Bek",
					Email:     "altyn@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/user", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			CreateUser(tc.store).ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("got status %d, want %d", rec.Code, tc.wantStatus)
			}

			var response struct {
				Code    int             `json:"code"`
				Comment string          `json:"comment"`
				Data    json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("response is not valid JSON: %v", err)
			}
			if response.Code != tc.wantStatus {
				t.Fatalf("got response code %d, want %d", response.Code, tc.wantStatus)
			}
		})
	}
}
