package todos

import (
	"Altyn/projectOne/internal/api"
	"Altyn/projectOne/internal/auth"
	"Altyn/projectOne/internal/database"
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
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

type mockUserQueryResult struct {
	user *database.User
	err  error
}

var (
	registerMockDriverOnce sync.Once
	mockDriverMu           sync.Mutex
	mockQueryResult        mockUserQueryResult
)

type mockDriver struct{}

func (d *mockDriver) Open(string) (driver.Conn, error) {
	return &mockConn{}, nil
}

type mockConn struct{}

func (c *mockConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not implemented") }
func (c *mockConn) Close() error                        { return nil }
func (c *mockConn) Begin() (driver.Tx, error)           { return nil, errors.New("not implemented") }

func (c *mockConn) QueryContext(_ context.Context, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	mockDriverMu.Lock()
	res := mockQueryResult
	mockDriverMu.Unlock()

	if res.err != nil {
		return nil, res.err
	}
	if res.user == nil {
		return &mockUserRows{done: true}, nil
	}
	return &mockUserRows{user: res.user}, nil
}

type mockUserRows struct {
	user *database.User
	done bool
}

func (r *mockUserRows) Columns() []string {
	return []string{"id", "useraname", "first_name", "last_name", "email", "password", "created_at", "updated_at"}
}

func (r *mockUserRows) Close() error { return nil }

func (r *mockUserRows) Next(dest []driver.Value) error {
	if r.done || r.user == nil {
		return io.EOF
	}
	r.done = true
	dest[0] = r.user.ID.String()
	dest[1] = r.user.Useraname
	dest[2] = r.user.FirstName
	dest[3] = r.user.LastName
	dest[4] = r.user.Email
	dest[5] = r.user.Password
	dest[6] = r.user.CreatedAt
	dest[7] = r.user.UpdatedAt
	return nil
}

func newMockQueries(t *testing.T, res mockUserQueryResult) *database.Queries {
	t.Helper()
	registerMockDriverOnce.Do(func() {
		sql.Register("mockuserdb", &mockDriver{})
	})

	mockDriverMu.Lock()
	mockQueryResult = res
	mockDriverMu.Unlock()

	db, err := sql.Open("mockuserdb", "")
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return database.New(db)
}

func TestLoginUserHandler(t *testing.T) {
	hashedPassword, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	testUser := &database.User{
		ID:        uuid.New(),
		Useraname: "altyn",
		FirstName: "Altyn",
		LastName:  "Bek",
		Email:     "altyn@example.com",
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name       string
		body       string
		dbResult   mockUserQueryResult
		wantStatus int
		wantToken  bool
	}{
		{
			name:       "malformed JSON",
			body:       "{",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email",
			body:       `{"email":"invalid-email","password":"password123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "password too short",
			body:       `{"email":"altyn@example.com","password":"short"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "user not found",
			body:       `{"email":"missing@example.com","password":"password123"}`,
			dbResult:   mockUserQueryResult{user: nil},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "database error",
			body:       `{"email":"altyn@example.com","password":"password123"}`,
			dbResult:   mockUserQueryResult{err: errors.New("db connection failed")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "wrong password",
			body:       `{"email":"altyn@example.com","password":"wrongpassword"}`,
			dbResult:   mockUserQueryResult{user: testUser},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "successful login with trimmed email",
			body:       `{"email":"  altyn@example.com  ","password":"password123"}`,
			dbResult:   mockUserQueryResult{user: testUser},
			wantStatus: http.StatusOK,
			wantToken:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &api.ApiConfig{
				Db:        newMockQueries(t, tc.dbResult),
				SecretKey: "test-secret-key",
			}

			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			LoginUserHandler(cfg).ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("got status %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
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

			if tc.wantToken {
				var userResp UserResponse
				if err := json.Unmarshal(response.Data, &userResp); err != nil {
					t.Fatalf("failed to unmarshal user response: %v", err)
				}
				if userResp.Token == "" {
					t.Fatal("expected non-empty JWT token")
				}
				if userResp.Email != testUser.Email {
					t.Fatalf("got email %q, want %q", userResp.Email, testUser.Email)
				}
			}
		})
	}
}
