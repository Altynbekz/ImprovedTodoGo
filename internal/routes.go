package internal

import (
	"Altyn/projectOne/internal/api"
	"Altyn/projectOne/internal/todos"
	"log"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewRouter(cfg *api.ApiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api", http.StripPrefix("/api", http.FileServer(http.Dir("."))))

	mux.Handle("POST /api/todo", todos.CreateTODO(cfg))
	mux.Handle("GET /api/todo", todos.GetAllTodo(cfg))
	log.Println("The /api/todo endpoints are ready to use")

	mux.Handle("POST /api/user", todos.CreateUser(cfg.Db))
	mux.Handle("DELETE /api/user", todos.DeleteUser(cfg))
	log.Println("The /api/user endpoints are ready to use")

	mux.Handle("/api/login", todos.LoginUserHandler(cfg))
	log.Println("Login endpoint is ready to use")

	return mux
}
