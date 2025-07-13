// Package util
// contains general utility functionality
package util

import (
	"regexp"
	"strings"
)

const SEP = " - "

// GenerateTrackName -> "artist1 & artist2 & ... - track"
func GenerateTrackName(artists []string, track string) string {
	return generateConcatNames(artists) + SEP + CleanName(track)
}

// GenerateAlbumName -> "artist1 & artist2 & ... - album"
func GenerateAlbumName(artists []string, album string) string {
	return generateConcatNames(artists) + SEP + CleanName(album)
}

func GenerateArtistName(artist string) string {
	return CleanName(artist)
}

// GenerateConcatNames []string{artist1, artist2, ...} -> "artist1 & artist2 & ..."
func generateConcatNames(names []string) string {
	if len(names) > 0 {
		return strings.Join(cleanNames(names), " & ")
	}
	return ""
}

// JSONArrToStrArr converts "["str1", "str2", ...]" -> []string{"str1", "str2", ...}
func JSONArrToStrArr(jsonArr string) []string {
	return strings.Split(strings.ReplaceAll(strings.Trim(jsonArr, "[]"), "\"", ""), ",")
}

func LogProgress(current int, total int) {
	Log.Debugf("completed(%d%%)\t%d/%d\n", int(float32(current)/float32(total)*100), current, total)
}

// only used once but more convienient here
func cleanNames(names []string) []string {
	for i, name := range names {
		names[i] = CleanName(name)
	}

	return names
}

func CleanName(name string) string {
	re := regexp.MustCompile(`\s*[-,&]\s*`)
	return re.ReplaceAllString(name, " ")
}
