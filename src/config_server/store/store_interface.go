package store

type Store interface {
	Put(key string, value string) error
	GetByName(name string) (Configuration, error)
	GetByID(id string) (Configuration, error)
	Delete(key string) (bool, error)
}
