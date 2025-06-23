package main

import (
	"context"
	"database/sql"
	"fmt"
	"main/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const MediaPath string = "../music/"

func main() {
	fmt.Println("Starting...")
	conn, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connecting to database...")
	queries := db.New(conn)
	fmt.Println("Connected to database")

	// engine := ingestion.NewEngine("../music/data", queries)
	// engine.IngestAndCreateData()

	fmt.Println("Getting all track items...")
	trackItems, err := queries.GetAllTrackItems(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(trackItems)

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/track-items", func(c echo.Context) error {
		return c.JSON(http.StatusOK, trackItems)
	})

	e.GET("/tracks/:id/image", serveTrackCoverImage(queries))
	e.GET("/tracks/:id/audio", serveTrackAudio(queries))
	e.GET("/album/:id/image", serveAlbumImage(queries))
	e.GET("/artists/:id/image", serveArtistImage(queries))

	e.Logger.Fatal(e.Start(":8080"))
}
