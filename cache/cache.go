package cache

type Cacher interface {
	GetCache(key string) ([]byte, error)
	SetCache(key string, value []byte) error
}
