package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jlauser/gocache/internal/config"
	"github.com/jlauser/gocache/internal/db"
	"net/http"
	"time"
)

type Server struct {
	Config *config.Config
	CSV    *db.CsvDB
	Cache  *db.MemoryDB
}

func (s *Server) NewServer() http.Handler {
	mux := chi.NewRouter()
	// middleware settings
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Timeout(time.Duration(s.Config.Api.IdleTimeout) * time.Second))
	// CORS settings
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	s.addRoutes(mux) // all routes are in the routes.go file to keep code clean
	return mux
}

func (s *Server) cacheGet(key string) (string, bool) {
	data, ok := s.Cache.Read(key)
	if ok {
		return data.(string), ok
	}
	return "", false
}

func (s *Server) cacheInsert(key string, value string) bool {
	err := s.Cache.Create(key, value)
	if err != nil {
		return false
	}
	return true
}

//
//func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(status)
//	if err := json.NewEncoder(w).Encode(v); err != nil {
//		return fmt.Errorf("error encoding json: %w", err)
//	}
//	return nil
//}
//
//func decode[T any](r *http.Request) (T, error) {
//	var v T
//	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
//		return v, fmt.Errorf("error decoding json: %w", err)
//	}
//	return v, nil
//}
