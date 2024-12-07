package db

import (
	"context"
)

type MemoryStore map[string]any

type MemoryDB struct {
	Ctx   context.Context
	Store MemoryStore
}

func InitializeMemoryDB() (*MemoryDB, error) {
	result := &MemoryDB{
		Ctx:   context.Background(),
		Store: make(MemoryStore),
	}
	return result, nil
}

func (db *MemoryDB) Create(key string, value any) error {
	db.Store[key] = value
	return nil
}

func (db *MemoryDB) Read(key string) (any, bool) {
	value, ok := db.Store[key]
	if !ok {
		return "", false
	}
	return value, true
}

func (db *MemoryDB) Find(key string) (any, bool) {
	value, ok := db.Store[key]
	if !ok {
		return "", false
	}
	return value, true
}

func (db *MemoryDB) Update(key string, value any) {
	db.Store[key] = value
}

func (db *MemoryDB) Delete(key string) {
	delete(db.Store, key)
}
