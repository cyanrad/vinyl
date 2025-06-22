package main

import (
	"context"
	"fmt"
	"log"
	"main/db"
	"main/ingestion"
	"main/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func serveTrackCoverImage(db *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		// getting track ID
		idStr := c.Param("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Track id "+idStr+" is not valid")
		}

		// getting track item
		trackItem, err := db.GetOneTrackItems(context.Background(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Track id "+idStr+" cannot be found")
		}

		path, exists := getTrackCoverFilePath(trackItem)
		if !exists {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}

		// Log the request
		log.Printf("Serving track cover: %s to %s", path, c.RealIP())

		// Set appropriate headers
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")

		// Serve the file
		return c.File(path)
	}
}

func getTrackCoverFilePath(trackItem db.GetOneTrackItemsRow) (string, bool) {
	artistName := strings.Split(trackItem.ArtistNames, ",")[0]

	// we attempt to find a cover for the track first, then ablum cover then artist image
	// track cover
	if path, exists := checkIfImageExists(ingestion.TRACKS,
		util.GenerateTrackName(trackItem.Title, artistName, trackItem.AlbumName)); exists {
		return path, true
	}

	// album cover
	if trackItem.AlbumName != nil {
		if path, exists := checkIfImageExists(ingestion.ALBUMS,
			util.GenerateAlbumName(artistName, *trackItem.AlbumName)); exists {
			return path, true
		}
	}

	// artist Image
	if path, exists := checkIfImageExists(ingestion.ARTISTS, artistName); exists {
		return path, true
	}

	return "", false
}

func checkIfImageExists(resourceType ingestion.IngestionType, resourceName string) (string, bool) {
	for _, extention := range []string{".png", ".jpg"} {
		// Construct file path
		filePath := filepath.Join("../music", string(resourceType), resourceName+extention)
		fmt.Println(filePath)

		// Check if file exists
		_, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		return filePath, true
	}

	return "", false
}
