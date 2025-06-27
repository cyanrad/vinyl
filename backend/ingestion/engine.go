// Package ingestion is meant for handling artist, tracks, albums, etc. data from json file
// and loading it into the database.
package ingestion

import (
	"encoding/json"
	"main/db"
	"main/util"
	"os"
	"path/filepath"
)

type Engine struct {
	cache   *Cache
	queries *db.Queries
}

func NewEngine(queries *db.Queries) *Engine {
	cache := NewCache()
	return &Engine{
		cache:   cache,
		queries: queries,
	}
}

func (e *Engine) IngestArtists() ([]ArtistIngestion, error) {
	dataPath := filepath.Join(util.DATA_PATH, string(util.ARTISTS)+".json")
	artistData, err := os.ReadFile(dataPath)
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
	dataPath := filepath.Join(util.DATA_PATH, string(util.TRACKS)+".json")
	trackData, err := os.ReadFile(dataPath)
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

func (e *Engine) IngestAlbums() ([]AlbumIngestion, error) {
	dataPath := filepath.Join(util.DATA_PATH, string(util.ALBUMS)+".json")
	albumData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var albums []AlbumIngestion
	err = json.Unmarshal(albumData, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
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

	albums, err := e.IngestAlbums()
	if err != nil {
		return err
	}

	err = e.CreateArtists(artists)
	if err != nil {
		return err
	}

	err = e.CreateAlbums(albums)
	if err != nil {
		return err
	}

	err = e.CreateTracks(tracks)
	if err != nil {
		return err
	}

	return nil
}
