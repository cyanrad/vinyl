package main

import (
	"context"
	"database/sql"
	"fmt"
	"main/db"
	"main/ingestion"

	_ "github.com/mattn/go-sqlite3"
)

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

	engine := ingestion.NewEngine("../music/data", queries)

	artists, err := engine.IngestArtists()
	if err != nil {
		panic(err)
	}

	tracks, err := engine.IngestTracks()
	if err != nil {
		panic(err)
	}

	fmt.Println(artists)
	engine.CreateArtists(artists)
	fmt.Println("Created artists")

	fmt.Println(tracks)
	err = engine.CreateTracks(tracks)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created tracks")

	db_tracks, err := queries.GetAllTracks(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(db_tracks)

	// Get all track items
	// fmt.Println("Getting all track items...")
	// trackItems, err := queries.GetAllTrackItems(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(trackItems)
}
