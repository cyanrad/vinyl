// Package spotify package for ingesting spotify data
package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/ingestion/storage"
	"main/util"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyConn struct {
	client *spotify.Client
	ctx    context.Context
	cache  *storage.Cache
}

func Connect(ctx context.Context, cache *storage.Cache) (*SpotifyConn, error) {
	if cache == nil {
		return nil, errors.New("nil cache passed to SpotifyConn.Connect")
	}

	// TODO: should handle token expiry
	if util.SPOTIFY_SECRET == "" || util.SPOTIFY_ID == "" {
		return nil, errors.New("spotify id or secret not provided")
	}

	// Credentials to get the API token
	config := &clientcredentials.Config{
		ClientID:     util.SPOTIFY_ID,
		ClientSecret: util.SPOTIFY_SECRET,
		TokenURL:     spotifyauth.TokenURL,
	}

	log.Println("Getting Spotify token")
	token, err := config.Token(ctx)
	if err != nil {
		return nil, errors.New("failed to get spotify token. check if api id & secret are valid")
	}

	log.Println("Creating Spotify client")
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	s := SpotifyConn{
		client: client,
		ctx:    ctx,
		cache:  cache,
	}
	return &s, nil
}

func (s *SpotifyConn) IngestPlaylist(playlistID string) error {
	data, err := s.GeneratePlaylistData(playlistID)
	log.Println(err)
	log.Println(data)

	artists, err := s.GetFullArtists(data.ArtistIDs)
	log.Println(err)
	for _, artist := range artists {
		log.Println(*artist)
	}

	albums, err := s.GetFullAlbums(data.AlbumIDs)
	log.Println(err)
	for _, album := range albums {
		log.Println(*album)
	}

	return err
}

func (s *SpotifyConn) GetTrackIngestion(trackID string) (TrackIngestion, error) {
	// getting the track data
	tracks, err := s.GetFullTracks([]spotify.ID{spotify.ID(trackID)})
	if err != nil {
		return TrackIngestion{}, err
	}
	if len(tracks) > 1 {
		panic("what??!?")
	}
	track := tracks[0]

	// getting artists who worked on the track
	trackArtists, err := s.generateSimpleArtistIngestion(track.Artists)
	if err != nil {
		return TrackIngestion{}, nil
	}

	// getting albums data
	var albumIngestion *AlbumIngestion = nil
	if track.Album.AlbumType != "single" { // if the album is a single we don't count it as an album in our data
		albums, err := s.GetFullAlbums([]spotify.ID{track.Album.ID})
		if err != nil {
			return TrackIngestion{}, err
		}
		if len(albums) > 1 {
			panic("whaaaaaaat??!??!?")
		}
		album := albums[0]

		// getting all the artists involved in the album
		albumArtists, err := s.generateSimpleArtistIngestion(album.Artists)
		if err != nil {
			return TrackIngestion{}, nil
		}

		// generating album ingestion
		albumIngestion = &AlbumIngestion{
			Name:     album.Name,
			ImageURL: album.Images[0].URL,
			Artists:  albumArtists,
		}
	}

	// creating imgestion type
	return TrackIngestion{
		Title:     track.Name,
		Tags:      nil,
		AlbumRank: int(track.TrackNumber),
		ImageURL:  track.Album.Images[0].URL, // getting the highest res image
		Album:     albumIngestion,
		Artists:   trackArtists,
	}, nil
}

// getCached the bool is supposed to indicate if data is usable
func (s *SpotifyConn) getCached(resource util.ResourceType, source util.IngestionSource, id string, data any) (bool, error) {
	if data == nil {
		return false, errors.New("nil cache passed to SpotifyConn.getCached")
	}

	log.Printf("Checking if resource %s from %s id %s is cached", resource, source, id)
	found, err := s.cache.Get(resource, source, id, data)
	if err != nil {
		if found {
			log.Printf("Unexpected error occoured even though object %s was found. Busting cache", id)
			delerr := s.cache.Delete(util.PLAYLISTS, util.SOURCE_SPOTIFY, id)
			if delerr != nil { // this is too nested...
				err = fmt.Errorf("%v; %v", err, delerr)
			}
		}

		return false, err
	}

	if found {
		log.Println("Cached object found")
	} else {
		log.Println("Cached object was not found")
	}
	return found, nil
}

// WARNING: please for the love of god find a better way to do this
func deduplicate(input []spotify.ID) []spotify.ID {
	seen := make(map[spotify.ID]bool)
	result := []spotify.ID{}

	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}
