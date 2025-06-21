package ingestion

type Cache struct {
	cache    map[string]interface{}
	dataPath string
}

func NewCache(dataPath string) *Cache {
	return &Cache{
		cache:    make(map[string]interface{}),
		dataPath: dataPath,
	}
}
