package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jlauser/gocache/api"
	"github.com/jlauser/gocache/internal/config"
	"github.com/jlauser/gocache/internal/db"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func panicHandler() {
	if r := recover(); r != nil {
		log.Println("recovered from error:", r)
	}
}

func run(ctx context.Context) error {
	var _, cancel = signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// load config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// prometheus metrics
	promMux := http.NewServeMux()
	promMux.Handle(cfg.Prometheus.Route, promhttp.Handler())
	promHttpServer := &http.Server{
		Addr:    cfg.Prometheus.Address,
		Handler: promMux,
	}
	go func() {
		log.Printf("Prometheus listening on %s", promHttpServer.Addr)
		if err := promHttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting prometheus server @ %s: %v\n", promHttpServer.Addr, err)
		}
	}()

	// CSV db
	csvDB, err := db.InitializeCsvDB(cfg.Data)
	if err != nil {
		panic(err)
	}

	cache, err := db.InitializeMemoryDB()
	if err != nil {
		panic(err)
	}

	// main API
	apiServer := &api.Server{
		Config: cfg,
		CSV:    csvDB,
		Cache:  cache,
	}
	apiMux := apiServer.NewServer()
	apiHttpServer := &http.Server{
		Addr:         cfg.Api.Address,
		Handler:      apiMux,
		WriteTimeout: time.Second * time.Duration(cfg.Api.WriteTimout),
		ReadTimeout:  time.Second * time.Duration(cfg.Api.ReadTimout),
		IdleTimeout:  time.Second * time.Duration(cfg.Api.IdleTimeout),
	}
	go func() {
		log.Printf("API listening on %s", apiHttpServer.Addr)
		if err := apiHttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("API failed to start @ %s: %v\n", apiHttpServer.Addr, err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Second*10)
		defer cancel()
		if err := promHttpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Prometheus failed to shutdown @ %s: %v\n", cfg.Prometheus.Address, err)
		}
		if err := apiHttpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("API failed to shutdown @ %s: %v\n", cfg.Api.Address, err)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	defer panicHandler()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%s\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
