package ingestion

import (
	"encoding/json"
	"main/db"
	"os"
)

type Engine struct {
	cache    *Cache
	queries  *db.Queries
	dataPath string
}

func NewEngine(dataPath string, queries *db.Queries) *Engine {
	cache := NewCache(dataPath)
	return &Engine{
		cache:    cache,
		queries:  queries,
		dataPath: dataPath,
	}
}

func (e *Engine) IngestArtists() ([]ArtistIngestion, error) {
	artistData, err := os.ReadFile(e.dataPath + "/" + string(INGESTION_TYPE_ARTISTS) + ".json")
	if err != nil {
		return nil, err
	}

	var artists []ArtistIngestion
	err = json.Unmarshal(artistData, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (e *Engine) IngestTracks() ([]TrackIngestion, error) {
	trackData, err := os.ReadFile(e.dataPath + "/" + string(INGESTION_TYPE_TRACKS) + ".json")
	if err != nil {
		return nil, err
	}

	var tracks []TrackIngestion
	err = json.Unmarshal(trackData, &tracks)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}
