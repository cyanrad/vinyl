// Package util
// contains general utility functionality
package util

import (
	"log"
	"strings"
)

const SEP = " - "

// GenerateTrackName -> "artist1 & artist2 & ... - album - track"
// to generate the correct name you should use GenerateConcatNames on artists before using this
// if album is nil it's omitted
func GenerateTrackName(artistName string, albumName *string, trackName string) string {
	name := artistName
	if albumName != nil {
		name += SEP + *albumName
	}
	name += SEP + trackName

	return name
}

// GenerateAlbumName -> "artist1 & artist2 & ... - album"
// to generate the correct name you should use GenerateConcatNames on artists before using this
func GenerateAlbumName(artistNames string, albumName string) string {
	return artistNames + SEP + albumName
}

// GenerateConcatNames []string{artist1, artist2, ...} -> "artist1 & artist2 & ..."
func GenerateConcatNames(names []string) string {
	return strings.Join(names, " & ")
}

// JSONArrToStrArr converts "["str1", "str2", ...]" -> []string{"str1", "str2", ...}
func JSONArrToStrArr(jsonArr string) []string {
	return strings.Split(strings.ReplaceAll(strings.Trim(jsonArr, "[]"), "\"", ""), ",")
}

func LogProgress(current int, total int) {
	log.Printf("completed(%d%%)\t%d/%d\n", int(float32(current)/float32(total)*100), current, total)
}
