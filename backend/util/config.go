package util

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	// API
	PORT int

	// FS
	MEDIA_PATH    string
	DATABASE_PATH string
	FRONTEND_PATH string
	DATA_PATH     string

	// Ingesetion
	INGEST      bool
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
	var relDBPath, relDataPath, relFrontendPath string
	flag.StringVar(&relDBPath, "database-path", "database.db", "SQLite database file path (relative to media-path)")
	flag.StringVar(&relDataPath, "data-path", "data", "Resource JSON data directory path (relative to media-path)")
	flag.StringVar(&relFrontendPath, "frontend-path", "dist", "Frontend dist path (relative to media-path)")

	// Ingestion and API usage should be seperate operations in the executable
	flag.BoolVar(&INGEST, "ingest", false, "Runs ingestion mode where the data is consumed and inserted into the database file specificed by the database-path arg")

	var source, resource string
	flag.StringVar(&source, "source", "", "The source the data is coming from. current available sources are: local, spotify")
	flag.StringVar(&resource, "resource", "", "The resource type being ingested (not needed for local).\nAvailable resource: tracks, artists, albums, playlists, audio")
	flag.StringVar(&RESOURCE_ID, "resource_id", "", "Resource ID to be ingested (not needed for local).")

	flag.Parse()

	// Convert relative paths to absolute paths based on MediaPath
	DATABASE_PATH = filepath.Join(MEDIA_PATH, relDBPath)
	DATA_PATH = filepath.Join(MEDIA_PATH, relDataPath)
	FRONTEND_PATH = filepath.Join(MEDIA_PATH, relFrontendPath)

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
