// Package spotify package for ingesting spotify data
package spotify

import (
	"context"
	"errors"
	"log"
	"main/util"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyConn struct {
	client *spotify.Client
	ctx    context.Context
}

func Connect(ctx context.Context) (SpotifyConn, error) {
	// TODO: should handle token expiry
	if util.SPOTIFY_SECRET == "" || util.SPOTIFY_ID == "" {
		return SpotifyConn{}, errors.New("spotify id or secret not provided")
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
		return SpotifyConn{}, errors.New("failed to get spotify token. check if api id & secret are valid")
	}

	log.Println("Creating Spotify client")
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return SpotifyConn{
		client: client,
		ctx:    ctx,
	}, nil
}
