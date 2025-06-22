// Package util
// contains general utility functionality
package util

const SEP = " - "

func GenerateTrackName(trackName string, artistName string, albumName *string) string {
	name := artistName
	if albumName != nil {
		name += SEP + *albumName
	}
	name += SEP + trackName

	return name
}

func GenerateAlbumName(artistName string, albumName string) string {
	return artistName + SEP + albumName
}
