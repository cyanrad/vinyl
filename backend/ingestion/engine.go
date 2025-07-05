// Package ingestion is meant for handling artist, tracks, albums, etc. data from json file
// and loading it into the database.
// TODO: currently this package is inconsistent between spotify and local, this should be fixed
package ingestion

import (
	"context"
	"encoding/json"
	"log"
	"main/db"
	"main/ingestion/spotify"
	"main/ingestion/storage"
	"main/util"
	"os"
	"path/filepath"
)

type Engine struct {
	cache   *storage.Cache
	queries *db.Queries
}

func NewEngine(queries *db.Queries, cachePath string) *Engine {
	cache := storage.NewCache(cachePath)
	return &Engine{
		cache:   cache,
		queries: queries,
	}
}

func (e *Engine) IngestSpotify(ctx context.Context, resourceType util.ResourceType, resourceID string) error {
	spotifyConn, err := spotify.Connect(ctx, e.cache)
	if err != nil {
		panic(err)
	}

	if util.BUST_CACHE {
		log.Printf("Removing cached object for Spotify resource: %s - ID: %s\n", string(resourceType), resourceID)
		err := e.cache.Delete(resourceType, util.SOURCE_SPOTIFY, resourceID)
		if err != nil {
			return err
		}
	}

	switch resourceType {
	case util.PLAYLISTS:
		data, err := spotifyConn.GeneratePlaylistData(resourceID)
		log.Println(err)
		log.Println(data)

		artists, err := spotifyConn.GetFullArtists(data.ArtistIDs)
		for _, artist := range artists {
			log.Println(*artist)
		}
		return err
	default:
		log.Fatalf("ingestion type %s for spotify is not implemented. terminating program", resourceType)
		return nil
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
