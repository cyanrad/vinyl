package ingestion

import "encoding/json"

type IngestionType string

const (
	INGESTION_TYPE_TRACKS    IngestionType = "tracks"
	INGESTION_TYPE_ARTISTS   IngestionType = "artists"
	INGESTION_TYPE_ALBUMS    IngestionType = "albums"
	INGESTION_TYPE_PLAYLISTS IngestionType = "playlists"
)

type ArtistIngestion struct {
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Links       ArtistLinks `json:"links"`
}

type ArtistLinks struct {
	SoundCloud *string `json:"soundcloud,omitempty"`
	Spotify    *string `json:"spotify,omitempty"`
	Personal   *string `json:"personal,omitempty"`
}

func (a ArtistLinks) ToString() string {
	json, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(json)
}

type TrackIngestion struct {
	Title       string   `json:"title"`
	Artists     []string `json:"artists"`
	Album       *string  `json:"album,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags"`
}
