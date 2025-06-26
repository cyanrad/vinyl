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

	"github.com/labstack/echo/v4"
)

func serveArtistImage(db *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		// getting artist ID
		idStr := c.Param("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Artist id "+idStr+" is not valid")
		}

		// getting artist data
		artist, err := db.GetArtistById(context.Background(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Artist id "+idStr+" cannot be found")
		}

		// getting artist image path
		path, exists := checkIfImageExists(ingestion.ARTISTS, artist.Name)
		if !exists {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}

		// Log the request
		log.Printf("Serving artist image: %s to %s", path, c.RealIP())

		// Set appropriate headers
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")

		// Serve the file
		return c.File(path)
	}
}

func serveAlbumImage(db *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		// getting album ID
		idStr := c.Param("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Album id "+idStr+" is not valid")
		}

		// getting album data
		album, err := db.GetAlbumById(context.Background(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Album id "+idStr+" cannot be found")
		}

		// getting album image path
		path, exists := checkIfImageExists(ingestion.ALBUMS,
			util.GenerateAlbumName(album.ArtistNames, album.Name))

		if !exists {
			return echo.NewHTTPError(http.StatusNotFound, "Album not found")
		}

		// Log the request
		log.Printf("Serving album image: %s to %s", path, c.RealIP())

		// Set appropriate headers
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")

		// Serve the file
		return c.File(path)
	}
}

func serveTrackCoverImage(db *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		// getting track ID
		idStr := c.Param("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Track id "+idStr+" is not valid")
		}

		// getting track item
		trackItem, err := db.GetTrackItemById(context.Background(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Track id "+idStr+" cannot be found")
		}

		// getting track cover file path
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

func serveTrackAudio(db *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		// getting track ID
		idStr := c.Param("id")

		// converting id to int
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Track id "+idStr+" is not valid")
		}

		// getting track item
		trackItem, err := db.GetTrackItemById(context.Background(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Track id "+idStr+" cannot be found")
		}

		// getting track cover file path
		artistNames := util.GenerateConcatNames(util.JSONArrToStrArr(trackItem.ArtistNames))
		path, exists := checkIfAudioExists(util.GenerateTrackName(trackItem.Title, artistNames, trackItem.AlbumName))
		if !exists {
			return echo.NewHTTPError(http.StatusNotFound, "Audio not found")
		}

		// Log the request
		log.Printf("Serving track audio: %s to %s", path, c.RealIP())

		// Set appropriate headers
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")

		// Serve the file
		return c.File(path)
	}
}

func getTrackCoverFilePath(trackItem db.GetTrackItemByIdRow) (string, bool) {
	artistName := util.GenerateConcatNames(util.JSONArrToStrArr(trackItem.ArtistNames))

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
		filePath := filepath.Join(MediaPath, string(resourceType), resourceName+extention)
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

func checkIfAudioExists(resourceName string) (string, bool) {
	// Construct file path
	filePath := filepath.Join(MediaPath, "audio", resourceName+".mp3")
	fmt.Println(filePath)

	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}

	return filePath, true
}
