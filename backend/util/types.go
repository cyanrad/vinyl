package util

type ResourceType string

const (
	TRACKS    ResourceType = "tracks"
	ARTISTS   ResourceType = "artists"
	ALBUMS    ResourceType = "albums"
	PLAYLISTS ResourceType = "playlists"
	AUDIO     ResourceType = "audio"
)
