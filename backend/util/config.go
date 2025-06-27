package util

import (
	"flag"
	"path/filepath"
)

var (
	// API
	PORT int

	// FS
	MEDIA_PATH    string
	DATABASE_PATH string
	FRONTEND_PATH string
	DATA_PATH     string
)

func InitConfig() {
	// Define command line flags with default values
	flag.IntVar(&PORT, "port", 8080, "Server port")
	flag.StringVar(&MEDIA_PATH, "media-path", "./files", "Base media directory path")

	// getting paths relative to media folder
	var relDBPath, relDataPath, relFrontendPath string
	flag.StringVar(&relDBPath, "database-path", "database.db", "SQLite database file path (relative to media-path)")
	flag.StringVar(&relDataPath, "data-path", "data", "Resource JSON data directory path (relative to media-path)")
	flag.StringVar(&relFrontendPath, "frontend-path", "dist", "Frontend dist path (relative to media-path)")

	// Parse command line arguments
	flag.Parse()

	// Convert relative paths to absolute paths based on MediaPath
	DATABASE_PATH = filepath.Join(MEDIA_PATH, relDBPath)
	DATA_PATH = filepath.Join(MEDIA_PATH, relDataPath)
	FRONTEND_PATH = filepath.Join(MEDIA_PATH, relFrontendPath)
}
