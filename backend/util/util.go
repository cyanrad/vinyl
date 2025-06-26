// Package util
// contains general utility functionality
package util

import (
	"strings"
)

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

func GenerateConcatNames(names []string) string {
	return strings.Join(names, " & ")
}

// JSONArrToStrArr converts "["str1", "str2", ...]" -> []string{"str1", "str2", ...}
func JSONArrToStrArr(jsonArr string) []string {
	return strings.Split(strings.ReplaceAll(strings.Trim(jsonArr, "[]"), "\"", ""), ",")
}
