package spotify

import (
	"main/ingestion/storage"

	"github.com/zmb3/spotify/v2"
)

const SPOTIFY_ARTIST_URL_BASE = "https://open.spotify.com/artist/"

func generateArtistIngestion(artist spotify.FullArtist) storage.ArtistIngestion {
	spotifyURL := string(SPOTIFY_ARTIST_URL_BASE + artist.ID)

	return storage.ArtistIngestion{
		Name:        artist.Name,
		Description: nil,
		Links: storage.ArtistLinks{
			Spotify: &spotifyURL,
		},
	}
}
