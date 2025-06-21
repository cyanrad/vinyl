package ingestion

type IngestionCache struct {
	cache     map[string]interface{}
	MediaPath string
}

func NewIngestionCache(mediaPath string) *IngestionCache {
	return &IngestionCache{
		cache:     make(map[string]interface{}),
		MediaPath: mediaPath,
	}
}
