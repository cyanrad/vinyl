package util

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// global consts
const (
	// spotify API
	PLAYLIST_PAGE_SIZE      int    = 100
	ARTIST_PAGE_SIZE        int    = 50
	TRACK_PAGE_SIZE         int    = 50
	ALBUM_PAGE_SIZE         int    = 20
	SPOTIFY_ARTIST_URL_BASE string = "https://open.spotify.com/artist/"
)

var (
	// API
	PORT int

	// FS
	MEDIA_PATH    string
	DATABASE_PATH string
	FRONTEND_PATH string
	DATA_PATH     string
	CACHE_PATH    string

	// Ingesetion
	INGEST      bool
	BUST_CACHE  bool
	SOURCE      IngestionSource
	RESOURCE    ResourceType
	RESOURCE_ID string

	// API Keys
	SPOTIFY_ID     string
	SPOTIFY_SECRET string
)

func InitConfig() {
	// === Flag Configs (Public Info) ===
	// Define command line flags with default values
	flag.IntVar(&PORT, "port", 8080, "Server port")
	flag.StringVar(&MEDIA_PATH, "media-path", "./files", "Base media directory path")

	// getting paths relative to media folder
	var relDBPath, relDataPath, relFrontendPath, relCachePath string
	flag.StringVar(&relDBPath, "database-path", "database.db", "SQLite database file path (relative to media-path)")
	flag.StringVar(&relDataPath, "data-path", "data", "Resource JSON data directory path (relative to media-path)")
	flag.StringVar(&relFrontendPath, "frontend-path", "dist", "Frontend dist path (relative to media-path)")
	flag.StringVar(&relCachePath, "cache-path", "cache", "Cache file that contains in memory objects & json files of external source APIs (relative to media-path)")

	// Ingestion and API usage should be seperate operations in the executable
	flag.BoolVar(&INGEST, "ingest", false, "Runs ingestion mode where the data is consumed and inserted into the database file specificed by the database-path arg")
	flag.BoolVar(&BUST_CACHE, "bust-cache", false, "Removes the cache on a given resource. best used for resetting an external source data")

	var source, resource string
	flag.StringVar(&source, "source", "", "The source the data is coming from. current available sources are: local, spotify")
	flag.StringVar(&resource, "resource", "", "The resource type being ingested (not needed for local).\nAvailable resource: tracks, artists, albums, playlists, audio")
	flag.StringVar(&RESOURCE_ID, "resource_id", "", "Resource ID to be ingested (not needed for local).")

	flag.Parse()

	// Convert relative paths to absolute paths based on MediaPath
	DATABASE_PATH = filepath.Join(MEDIA_PATH, relDBPath)
	DATA_PATH = filepath.Join(MEDIA_PATH, relDataPath)
	FRONTEND_PATH = filepath.Join(MEDIA_PATH, relFrontendPath)
	CACHE_PATH = filepath.Join(MEDIA_PATH, relCachePath)

	// Handling ingestion variables & checks
	SOURCE = MapStrToIngestionSource(source)
	if INGEST && SOURCE == SOURCE_UNKNOWN {
		log.Fatalf("ingestion source %s unknown. terminating program\n", source)
	}

	RESOURCE = MapStrToResourceType(resource)
	if INGEST && SOURCE != SOURCE_LOCAL && RESOURCE == UNKNOWN {
		log.Fatalf("ingestion resource %s unknown. terminating program\n", resource)
	}

	if INGEST && SOURCE != SOURCE_LOCAL && SOURCE != SOURCE_UNKNOWN && RESOURCE_ID == "" {
		log.Fatalln("ingestion on an external source was activated but no id was given. terminating program")
	}

	// === Env Configs (Sensative Info) ===
	// Load .env file from MEDIA_PATH if it exists, ignore error if file doesn't exist
	err := godotenv.Load(filepath.Join(MEDIA_PATH, ".env"))
	if err != nil && SOURCE != SOURCE_LOCAL {
		log.Println("file/.env not found. API keys expected to be provided through env vars. Otherwise external ingestion fails")
	}

	// Load Spotify API credentials from environment variables
	SPOTIFY_ID = os.Getenv("SPOTIFY_ID")
	SPOTIFY_SECRET = os.Getenv("SPOTIFY_SECRET")
}
