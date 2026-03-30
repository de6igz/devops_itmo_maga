package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("data", "games.db")
	}

	app, err := newApp(dbPath, true)
	if err != nil {
		log.Fatalf("failed to build application: %v", err)
	}
	defer app.Close()

	log.Printf("backend started on http://localhost:%s", port)
	if err := app.Server.Start(":" + port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
