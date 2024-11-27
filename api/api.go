package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlauser/gocache/internal/config"
	"github.com/jlauser/gocache/internal/db"
	"log"
	"net/http"
	"time"
)

type Application struct {
	Config *config.Config
	CSV    *db.CsvDB
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(app.Config.Api.IdleTimeout) * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.HealthCheckHandler)
	})
	return r
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.Config.Api.Address,
		Handler:      mux,
		WriteTimeout: time.Second * time.Duration(app.Config.Api.WriteTimout),
		ReadTimeout:  time.Second * time.Duration(app.Config.Api.ReadTimout),
		IdleTimeout:  time.Second * time.Duration(app.Config.Api.IdleTimeout),
	}

	log.Printf("API listening on %s", app.Config.Api.Address)
	return srv.ListenAndServe()
}
