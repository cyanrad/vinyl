package spotify

import "main/ingestion/storage"

type TrackIngestion struct {
	Title     string
	Tags      []string
	AlbumRank int
	ImageURL  string
	Artists   []ArtistIngestion
	Album     *AlbumIngestion
}

type AlbumIngestion struct {
	Name     string
	Artists  []ArtistIngestion
	ImageURL string
}

type ArtistIngestion struct {
	Name     string
	Links    storage.ArtistLinks
	ImageURL string
}
