package spotify

import "main/ingestion/storage"

type TrackIngestion struct {
	Title     string
	FullTitle string
	Tags      []string
	AlbumRank int
	ImageURL  string
	Artists   []string
	Album     *string
}

type AlbumIngestion struct {
	Name     string
	FullName string
	Artists  []string
	ImageURL string
}

type ArtistIngestion struct {
	Name     string
	Links    storage.ArtistLinks
	ImageURL string
}

type SpotifyIngestion struct {
	Artists []ArtistIngestion
	Albums  []AlbumIngestion
	Tracks  []TrackIngestion
}
