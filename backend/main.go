package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/db"
	"main/ingestion"
	"main/util"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Loading environment configs like port and file paths
	util.InitConfig()

	fmt.Println("Starting...")
	fmt.Println("Connecting to database...")
	conn, err := sql.Open("sqlite3", util.DATABASE_PATH)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	queries := db.New(conn)
	fmt.Println("Connected to database")

	// TODO: make this a command line option
	engine := ingestion.NewEngine(queries)
	engine.IngestAndCreateData()

	fmt.Println("Getting all track items...")
	trackItems, err := queries.GetAllTrackItems(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	// TODO: this should be updated with the server's domain at some point
	e.Use(middleware.CORS())

	// serving svelte frontend build
	e.Static("/", util.FRONTEND_PATH)

	// Serving basically all the data.
	// WARNING: this can become a bottleneck at large datasets, but currently it's not an issue
	e.GET("/track-items", func(c echo.Context) error {
		return c.JSON(http.StatusOK, trackItems)
	})

	// serving the audio
	// TODO: there should be another WS path to stream the audio
	e.GET("/tracks/:id/audio", serveTrackAudio(queries))

	// serving resource images
	// TODO: serving lower res images can be useful to lower data usage & improve load times
	e.GET("/tracks/:id/image", serveTrackCoverImage(queries))
	e.GET("/albums/:id/image", serveAlbumImage(queries))
	e.GET("/artists/:id/image", serveArtistImage(queries))

	// starting server with logging
	fmt.Println("Starting Vinyl :)")
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(util.PORT)))
}
