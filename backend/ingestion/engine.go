// Package ingestion is meant for handling artist, tracks, albums, etc. data from json file
// and loading it into the database.
// TODO: currently this package is inconsistent between spotify and local, this should be fixed
package ingestion

import (
	"context"
	"encoding/json"
	"main/db"
	"main/ingestion/storage"
	"main/util"
	"os"
	"path/filepath"
)

type Engine struct {
	cache   *storage.Cache
	queries *db.Queries
	ctx     context.Context
}

func NewEngine(ctx context.Context, queries *db.Queries, cachePath string) *Engine {
	cache := storage.NewCache(cachePath)
	return &Engine{
		cache:   cache,
		queries: queries,
		ctx:     ctx,
	}
}

func (e *Engine) IngestArtists() ([]storage.ArtistIngestion, error) {
	dataPath := filepath.Join(util.DATA_PATH, string(util.ARTISTS)+".json")
	artistData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var artists []storage.ArtistIngestion
	err = json.Unmarshal(artistData, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (e *Engine) IngestTracks() ([]storage.TrackIngestion, error) {
	dataPath := filepath.Join(util.DATA_PATH, string(util.TRACKS)+".json")
	trackData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var tracks []storage.TrackIngestion
	err = json.Unmarshal(trackData, &tracks)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (e *Engine) IngestAlbums() ([]storage.AlbumIngestion, error) {
	dataPath := filepath.Join(util.DATA_PATH, string(util.ALBUMS)+".json")
	albumData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var albums []storage.AlbumIngestion
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
