package main

import (
	"context"
	"database/sql"
	"fmt"
	"main/db"
	"net/http"

	"github.com/labstack/echo/v4"
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

	// engine := ingestion.NewEngine("../music/data", queries)

	fmt.Println("Getting all track items...")
	trackItems, err := queries.GetAllTrackItems(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(trackItems)

	e := echo.New()
	e.GET("/track-item", func(c echo.Context) error {
		return c.JSON(http.StatusOK, trackItems)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
