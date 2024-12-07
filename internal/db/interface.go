package db

type DB interface {
	Create(key string, value any) bool
	Read(key string) (any, bool)
	Find(key string) (any, bool)
	Update(key string, value any) bool
	Delete(key string) bool
}
