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
	fmt.Println("Starting...")

	// Loading environment configs like port and file paths
	util.InitConfig()

	fmt.Println("Connecting to database...")
	database := initDatabase()
	defer database.Conn.Close()
	fmt.Println("Connected to database")

	ctx := context.Background()
	if util.INGEST {
		runIngestion(ctx, database.Queries)
	} else {
		runServer(ctx, database.Queries)
	}
}

type databaseConn struct {
	Conn    *sql.DB
	Queries *db.Queries
}

func initDatabase() databaseConn {
	conn, err := sql.Open("sqlite3", util.DATABASE_PATH)
	if err != nil {
		panic(err)
	}

	queries := db.New(conn)
	return databaseConn{
		Conn:    conn,
		Queries: queries,
	}
}

func runIngestion(ctx context.Context, queries *db.Queries) {
	engine := ingestion.NewEngine(queries)

	switch util.SOURCE {
	case util.SOURCE_LOCAL:
		log.Println("Ingesting local data")
		engine.IngestAndCreateData()
	case util.SOURCE_SPOTIFY:
		log.Println("Ingesting Spotify data")
		engine.IngestSpotify(ctx, util.RESOURCE, util.RESOURCE_ID)
	default:
		panic("what the actual fuck")
	}
}

func runServer(ctx context.Context, queries *db.Queries) {
	e := echo.New()

	// TODO: this should be updated with the server's domain at some point
	e.Use(middleware.CORS())

	// serving svelte frontend build
	e.Static("/", util.FRONTEND_PATH)

	// Serving basically all the data.
	// WARNING: this can become a bottleneck at large datasets, but currently it's not an issue
	e.GET("/track-items", func(c echo.Context) error {
		log.Println("Getting all track items...")
		trackItems, err := queries.GetAllTrackItems(ctx)
		if err != nil {
			log.Fatal(err)
		}
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
