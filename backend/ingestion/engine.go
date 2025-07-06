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

func (e *Engine) IngestSpotify(resourceType util.ResourceType, resourceID string) error {
	spotifyConn, err := spotify.Connect(e.ctx, e.cache)
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

	log.Printf("Beginning ingestion of spotify %s %s", resourceType, resourceID)
	switch resourceType {
	case util.PLAYLISTS:
		err = spotifyConn.IngestPlaylist(resourceID)
	case util.TRACKS:
		err = e.IngestSpotifyTrack(spotifyConn, resourceID)
	default:
		log.Fatalf("ingestion type %s for spotify is not implemented. terminating program", resourceType)
	}

	return err
}

func (e *Engine) IngestSpotifyTrack(s *spotify.SpotifyConn, trackID string) error {
	ingestion, err := s.GetTrackIngestion(trackID)
	if err != nil {
		return err
	}

	// creating track artists
	trackArtists := make([]string, len(ingestion.Artists))
	for i, artist := range ingestion.Artists {
		trackArtists[i] = artist.Name
		if err = e.CreateArtists([]storage.ArtistIngestion{
			{Name: artist.Name, Links: artist.Links, Description: nil},
		}); err != nil {
			return err
		}

		storage.FetchImage(util.ARTISTS, artist.ImageURL, artist.Name)
	}

	// creating album
	var album *string = nil
	var albumRank *int = nil
	if ingestion.Album != nil {
		// creating album artists & getting their names
		artists := make([]string, len(ingestion.Album.Artists))
		for i, artist := range ingestion.Album.Artists {
			artists[i] = artist.Name
			if err = e.CreateArtists([]storage.ArtistIngestion{
				{Name: artist.Name, Links: artist.Links, Description: nil},
			}); err != nil {
				return err
			}
		}

		if err := e.CreateAlbums([]storage.AlbumIngestion{
			{Name: ingestion.Album.Name, Artists: artists, Description: nil},
		}); err != nil {
			return err
		}

		albumFullName := util.GenerateAlbumName(util.GenerateConcatNames(artists), ingestion.Album.Name)
		album = &albumFullName
		albumRank = &ingestion.AlbumRank

		// fetching album image
		storage.FetchImage(util.ALBUMS, ingestion.Album.ImageURL, albumFullName)
	}

	// creating track
	if err := e.CreateTracks([]storage.TrackIngestion{
		{
			Title:       ingestion.Title,
			Artists:     trackArtists,
			Album:       album,
			AlbumRank:   albumRank,
			Description: nil,
			Tags:        nil,
		},
	}); err != nil {
		return err
	}

	// we fetch the image if the album doesn't exist or is a single
	if ingestion.Album == nil {
		trackFullTitle := util.GenerateAlbumName(util.GenerateConcatNames(trackArtists), ingestion.Title)
		storage.FetchImage(util.TRACKS, ingestion.ImageURL, trackFullTitle)
	}

	return nil
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
