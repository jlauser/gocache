package db

type DB interface {
	Create(key string, value interface{}) bool
	Read(key string) (interface{}, bool)
	Find(key string) (interface{}, bool)
	Update(key string, value interface{}) bool
	Delete(key string) bool
}
