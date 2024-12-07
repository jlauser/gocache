package api

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) addRoutes(mux *chi.Mux) {
	// routes
	mux.Route("/v1", func(r chi.Router) {
		//r.Get("/health", s.HealthCheckHandler)
		r.Get("/users", s.UsersGetHandler)
		r.Get("/users/{id}", s.UsersGetHandler)
		r.Get("/search", s.SearchHandler)
		r.Get("/search/{q}", s.SearchHandler)
	})
}
