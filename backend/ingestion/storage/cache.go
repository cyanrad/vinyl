package storage

type Cache struct {
	cache map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]interface{}),
	}
}
