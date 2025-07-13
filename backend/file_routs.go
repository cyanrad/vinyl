package main

import (
	"context"
	"log"
	"main/db"
	"main/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

type api struct {
	db *db.Queries
}

func createAPI(db *db.Queries) api {
	return api{db}
}

type mediaPath struct {
	name     string
	resource util.ResourceType
}

func (a *api) serveResourceMedia(c echo.Context) error {
	// getting resource type
	resource := util.MapStrToResourceType(c.Param("resource"))
	if resource == util.UNKNOWN {
		return errorUnknownResource(c.Param("resource"))
	}

	// getting media type
	mediaType := util.MapStrToMediaType(c.Param("media"))
	if mediaType == util.MEDIA_UNKNOWN {
		return errorUnknownMediaType(c.Param("media"))
	}

	// audio is limited to tracks
	if mediaType == util.MEDIA_AUDIO && resource != util.TRACKS {
		return errorUnavailableMediaType(mediaType, resource)
	}

	// getting resource ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errorInvalidRsrouceID(resource, id)
	}

	log.Printf("INFO: getting media %s %s id '%d'", resource, mediaType, id)

	// Getting the resource name to fetch the file path
	// TODO: Implement for playlist
	names := []mediaPath{} // array to handle a unique edge case for track images - could be used for aliasing later
	switch resource {
	case util.ARTISTS:
		var artist db.Artist // need to do this so we don't repeat the error logic
		artist, err = a.db.GetArtistById(context.Background(), id)
		names = append(names, mediaPath{artist.Name, util.ARTISTS})
	case util.ALBUMS:
		var album db.GetAlbumByIdRow
		album, err = a.db.GetAlbumById(context.Background(), id)
		names = append(names, mediaPath{album.FullName, util.ALBUMS})
	case util.TRACKS:
		var track db.GetTrackItemByIdRow
		track, err = a.db.GetTrackItemById(context.Background(), id)
		names = getTrackNames(track)
	default:
		return errorUnknownResource(c.Param("resource"))
	}
	if err != nil {
		return errorResourceNotFound(resource, id)
	}

	log.Printf("INFO: got resource '%s'", names[0].name)

	// getting the media file path
	path := ""
	exists := false
	switch mediaType {
	case util.MEDIA_IMAGE:
		path, exists = getImagePath(names)
	case util.MEDIA_AUDIO:
		path, exists = getAudioPath(names[0])
	}
	if !exists {
		return errorMediaNotFound(mediaType, names[0])
	}

	log.Printf("INFO: sending file at '%s'", path)

	// Set appropriate headers
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")

	// Serve the file
	return c.File(path)
}

func getImagePath(paths []mediaPath) (string, bool) {
	for _, path := range paths {
		for _, extention := range []string{".png", ".jpg"} {
			// Construct file path
			filePath := filepath.Join(util.MEDIA_PATH, string(path.resource), path.name+extention)

			// Check if file exists
			_, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			return filePath, true
		}
	}
	return "", false
}

func getAudioPath(path mediaPath) (string, bool) {
	// Construct file path
	filePath := filepath.Join(util.MEDIA_PATH, string(util.AUDIO), path.name+".mp3")

	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}

	return filePath, true
}

func getTrackNames(track db.GetTrackItemByIdRow) []mediaPath {
	// I know it's not the best way to handle this but fuck it
	// it is important that this is ordered like track -> album -> artist as it sets fetch priority
	names := []mediaPath{{track.FullTitle, util.TRACKS}}
	if track.AlbumFullName != nil {
		names = append(names, mediaPath{*track.AlbumFullName, util.ALBUMS})
	}
	for _, artist := range util.JSONArrToStrArr(track.ArtistNames) {
		names = append(names, mediaPath{artist, util.ARTISTS})
	}
	return names
}

func errorInvalidRsrouceID(resource util.ResourceType, id int64) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, string(resource)+" id '"+string(id)+"' is not valid")
}

func errorResourceNotFound(resource util.ResourceType, id int64) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, string(resource)+" id '"+string(id)+"' cannot be found")
}

func errorUnknownResource(resource string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, "Unknown resource given '"+string(resource)+"'")
}

func errorUnknownMediaType(mediaType string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, "Unknown media type given '"+string(mediaType)+"'")
}

func errorUnavailableMediaType(mediaType util.MediaType, resource util.ResourceType) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, "Unavailable media type '"+string(mediaType)+"' for '"+string(resource)+"'")
}

func errorMediaNotFound(mediaType util.MediaType, path mediaPath) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, "Media "+string(mediaType)+" not found for "+string(path.resource)+" '"+path.name+"'")
}
