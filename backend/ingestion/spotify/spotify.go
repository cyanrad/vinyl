// Package spotify package for ingesting spotify data
package spotify

import (
	"context"
	"errors"
	"fmt"
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

	util.Log.Info("Getting Spotify token")
	token, err := config.Token(ctx)
	if err != nil {
		return nil, errors.New("failed to get spotify token. check if api id & secret are valid")
	}

	util.Log.Info("Creating Spotify client")
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	s := SpotifyConn{
		client: client,
		ctx:    ctx,
		cache:  cache,
	}
	return &s, nil
}

// getCached the bool is supposed to indicate if data is usable
func (s *SpotifyConn) getCached(resource util.ResourceType, source util.IngestionSource, id string, data any) (bool, error) {
	if data == nil {
		return false, errors.New("nil cache passed to SpotifyConn.getCached")
	}

	util.Log.Infof("Checking if resource %s from %s id %s is cached", resource, source, id)
	found, err := s.cache.Get(resource, source, id, data)
	if err != nil {
		if found {
			util.Log.Errorf("Unexpected error occoured even though object %s was found. Busting cache", id)
			delerr := s.cache.Delete(util.PLAYLISTS, util.SOURCE_SPOTIFY, id)
			if delerr != nil { // this is too nested...
				err = fmt.Errorf("%v; %v", err, delerr)
			}
		}

		return false, err
	}

	if found {
		util.Log.Debugf("Cache found %s %s '%s' ", source, resource, id)
	} else {
		util.Log.Debugf("Cache not found %s %s '%s'", source, resource, id)
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

// needed due to conflicting package names
func StrToID(id string) []spotify.ID {
	return []spotify.ID{spotify.ID(id)}
}
