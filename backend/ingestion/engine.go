package ingestion

import "database/sql"

type IngestionType string

const (
	INGESTION_TYPE_TRACK    IngestionType = "track"
	INGESTION_TYPE_ARTIST   IngestionType = "artist"
	INGESTION_TYPE_ALBUM    IngestionType = "album"
	INGESTION_TYPE_PLAYLIST IngestionType = "playlist"
)

type IngestionEngine struct {
	cache     *IngestionCache
	db        *sql.DB
	MediaPath string
}

func NewIngestionEngine(mediaPath string, db *sql.DB) *IngestionEngine {
	cache := NewIngestionCache(mediaPath)
	return &IngestionEngine{
		cache:     cache,
		db:        db,
		MediaPath: mediaPath,
	}
}

func (e *IngestionEngine) IngestAllArtists() []error {
	errs := []error{}

	return errs
}
