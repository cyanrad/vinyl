package util

import (
	"strings"
)

type ResourceType string

const (
	TRACKS    ResourceType = "tracks"
	ARTISTS   ResourceType = "artists"
	ALBUMS    ResourceType = "albums"
	PLAYLISTS ResourceType = "playlists"
	AUDIO     ResourceType = "audio"
	UNKNOWN   ResourceType = "unknown"
)

var resourceTypeMap = map[string]ResourceType{
	"track":    TRACKS,
	"artist":   ARTISTS,
	"album":    ALBUMS,
	"playlist": PLAYLISTS,
	"audio":    AUDIO,
}

func MapStrToResourceType(str string) ResourceType {
	if rt, ok := resourceTypeMap[str]; ok {
		return rt
	}
	return UNKNOWN
}

type IngestionSource string

const (
	SOURCE_UNKNOWN IngestionSource = "unknown"
	SOURCE_LOCAL   IngestionSource = "local"
	SOURCE_SPOTIFY IngestionSource = "spotify"
)

func MapStrToIngestionSource(str string) IngestionSource {
	str = strings.ToLower(str)
	switch str {
	case "local":
		return SOURCE_LOCAL
	case "spotify":
		return SOURCE_SPOTIFY
	default:
		return SOURCE_UNKNOWN
	}
}

type MediaType string

const (
	MEDIA_IMAGE   MediaType = "image"
	MEDIA_AUDIO   MediaType = "audio"
	MEDIA_UNKNOWN MediaType = "unknown"
)

func MapStrToMediaType(str string) MediaType {
	str = strings.ToLower(str)
	switch str {
	case "image":
		return MEDIA_IMAGE
	case "audio":
		return MEDIA_AUDIO
	default:
		return MEDIA_UNKNOWN
	}
}
