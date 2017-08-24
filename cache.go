package cache

import "time"

type Cache interface {
	Get(key string) interface{}
	Put(key string, val interface{}, timeout time.Duration) error
	Delete(key string) error
	IsExist(key string) bool
}
