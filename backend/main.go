package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/game_catalog?sslmode=disable"
	}

	app, err := newApp(databaseURL, true)
	if err != nil {
		log.Fatalf("failed to build application: %v", err)
	}
	defer app.Close()

	log.Printf("backend started on http://localhost:%s", port)
	if err := app.Server.Start(":" + port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
