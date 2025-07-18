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

	util.Log.Info("Connecting to database...")
	database := initDatabase()
	defer database.Conn.Close()
	util.Log.Info("Connected to database")

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
	engine := ingestion.NewEngine(ctx, queries, util.CACHE_PATH)

	switch util.SOURCE {
	case util.SOURCE_LOCAL:
		util.Log.Debug("Ingesting local data")
		engine.IngestAndCreateData()
	case util.SOURCE_SPOTIFY:
		util.Log.Debug("Ingesting Spotify data")
		err := engine.IngestSpotify(util.RESOURCE, util.RESOURCE_ID)
		if err != nil {
			util.Log.Errorf("Ingestion failed with: %s", err)
		}
	default:
		panic("what the actual fuck")
	}
}

func runServer(ctx context.Context, queries *db.Queries) {
	e := echo.New()
	a := createAPI(queries)

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

	// serving media (image, audio)
	e.GET("/:resource/:id/:media", a.serveResourceMedia)

	// starting server with logging
	fmt.Println("Starting Vinyl :)")
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(util.PORT)))
}
