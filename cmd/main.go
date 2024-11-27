package main

import (
	//"encoding/json"
	//"fmt"

	"errors"
	"github.com/jlauser/gocache/api"
	"github.com/jlauser/gocache/internal/config"
	"github.com/jlauser/gocache/internal/db"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	//"net/http"
)

func panicHandler() {
	if r := recover(); r != nil {
		log.Println("recovered from error:", r)
	}
}

func main() {
	defer panicHandler()

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// prometheus metrics
	promMux := http.NewServeMux()
	promMux.Handle(cfg.Prometheus.Route, promhttp.Handler())
	promServer := &http.Server{
		Addr:    cfg.Prometheus.Address,
		Handler: promMux,
	}

	// CSV db
	csvDB, err := db.InitializeCsvDB(cfg.Data)
	if err != nil {
		panic(err)
	}

	//key, ok := csvDB.Create("log:", []string{"test 1"})
	//if ok {
	//	result, ok := csvDB.Read("log:" + key)
	//	if ok {
	//		data := result.([]string)
	//		data[1] = "updated"
	//		csvDB.Update("log:"+key, data)
	//	}
	//	csvDB.Delete("log:" + key)
	//}

	cache, err := db.InitializeMemoryDB()
	if err != nil {
		panic(err)
	}

	// main API
	apiApp := &api.Application{
		Config: cfg,
		CSV:    csvDB,
		Cache:  cache,
	}
	apiMux := apiApp.Mount()

	// Start servers in separate go routine
	go func() {
		log.Printf("Prometheus listening on %s", promServer.Addr)
		if err := promServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting prometheus server @ %s: %v\n", promServer.Addr, err)
		}
	}()

	go func() {
		if err := apiApp.Run(apiMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting API server @ %s: %v\n", apiApp.Config.Api.Address, err)
		}
	}()
	select {}
}

/*
func startHttp() {
	// prometheus metrics
	promMux := http.NewServeMux()
	promMux.Handle("/metrics", promhttp.Handler())
	promServer := &http.Server{
		Addr:    ":2112",
		Handler: promMux,
	}

	// Set up second server with handlerTwo on port 9090
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", handleRoot)
	apiMux.HandleFunc("/query", handleQuery)
	apiMux.HandleFunc("/create", handleCreate)
	apiMux.HandleFunc("/read", handleRead)
	apiMux.HandleFunc("/update", handleUpdate)
	apiMux.HandleFunc("/delete", handleDelete)
	apiServer := &http.Server{
		Addr:    ":8080",
		Handler: apiMux,
	}

	// Start servers in separate go routine
	go func() {
		if err := promServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error starting prometheus server on port 2112: %v\n", err)
		}
	}()
	go func() {
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error starting API server on port 8080: %v\n", err)
		}
	}()

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the in-memory db store server!")
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Query")
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var db map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := db["key"].(string)
	value := db["value"]
	store[key] = value

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Key-value pair created successfully!")
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	value, ok := store[key]
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Value: %v", value)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	var db map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := db["key"].(string)
	value := db["value"]
	store[key] = value

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Key-value pair updated successfully!")
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	delete(store, key)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Key-value pair deleted successfully!")
}

*/
