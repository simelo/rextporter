package core

type RextKeyValueStore interface {
	GetString(key string) (string, error)
	SetString(key string, value string) (bool, error)
	GetObject(key string) (string, error)
	SetObject(key string, value interface{}) (bool, error)
}

type RextAuth interface{}
