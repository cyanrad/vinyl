// Package ingestion is meant for handling artist, tracks, albums, etc. data from json file
// and loading it into the database.
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
	artistData, err := os.ReadFile(e.dataPath + "/" + string(ARTISTS) + ".json")
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
	trackData, err := os.ReadFile(e.dataPath + "/" + string(TRACKS) + ".json")
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

func (e *Engine) IngestAndCreateData() error {
	artists, err := e.IngestArtists()
	if err != nil {
		return err
	}

	tracks, err := e.IngestTracks()
	if err != nil {
		return err
	}

	err = e.CreateArtists(artists)
	if err != nil {
		return err
	}

	err = e.CreateTracks(tracks)
	if err != nil {
		return err
	}

	return nil
}
