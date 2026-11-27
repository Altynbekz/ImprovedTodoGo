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
	log.Println(`Mapped "GET /api" -> StaticFileServer`)

	mux.Handle("POST /api/todo", todos.CreateTODO(cfg))
	log.Println(`Mapped "POST /api/todo" -> CreateTODO`)

	mux.Handle("GET /api/todo", todos.GetAllTodo(cfg))
	log.Println(`Mapped "GET /api/todo" -> GetAllTodo`)

	mux.Handle("POST /api/user", todos.CreateUser(cfg.Db))
	log.Println(`Mapped "POST /api/user" -> CreateUser`)

	mux.Handle("PUT /api/user", todos.UpdateUser(cfg.Db, cfg.SecretKey))
	log.Println(`Mapped "PUT /api/user" -> UpdateUser`)

	mux.Handle("DELETE /api/user", todos.DeleteUser(cfg))
	log.Println(`Mapped "DELETE /api/user" -> DeleteUser`)

	mux.Handle("POST /api/login", todos.LoginUserHandler(cfg))
	log.Println(`Mapped "POST /api/login" -> LoginUserHandler`)

	return mux
}
