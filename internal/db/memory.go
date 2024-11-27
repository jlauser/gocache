package db

import (
	"context"
)

type MemoryStore map[string]interface{}

var store = make(MemoryStore)

type MemoryDB struct {
	Ctx context.Context
}

func (db *MemoryDB) Create(key string, value interface{}) error {
	store[key] = value
	return nil
}

func (db *MemoryDB) Read(key string) (interface{}, bool) {
	value, ok := store[key]
	if !ok {
		return "", false
	}
	return value, true
}

func (db *MemoryDB) Find(key string) (interface{}, bool) {
	value, ok := store[key]
	if !ok {
		return "", false
	}
	return value, true
}

func (db *MemoryDB) Update(key string, value interface{}) {
	store[key] = value
}

func (db *MemoryDB) Delete(key string) {
	delete(store, key)
}

/*
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the in-memory db store server!")
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Query")
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := data["key"].(string)
	value := data["value"]
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
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := data["key"].(string)
	value := data["value"]
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
