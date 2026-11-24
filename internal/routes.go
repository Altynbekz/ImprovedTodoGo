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

	mux.Handle("/api/todo", todos.CreateTODO(cfg))
	log.Println("The /api/todo is ready to use")

	mux.Handle("/api/user", todos.CreateUser(cfg.Db))
	log.Println("The /api/user create endpoint is ready to use")

	return mux
}
