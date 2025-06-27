package ingestion

import "encoding/json"

type TrackIngestion struct {
	Title       string   `json:"title"`
	Artists     []string `json:"artists"`
	Album       *string  `json:"album,omitempty"`
	AlbumRank   *int     `json:"albumRank,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags"`
}

type AlbumIngestion struct {
	Name        string   `json:"name"`
	Artists     []string `json:"artists"`
	Description *string  `json:"description,omitempty"`
}

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
