package main

import (
	"context"
	"database/sql"
	"fmt"
	"main/db"

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

	// Get all track items
	fmt.Println("Getting all track items...")
	trackItems, err := queries.GetAllTrackItems(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(trackItems)
}
