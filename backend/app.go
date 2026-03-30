package main

import (
	"game-catalog-backend/internal/game"
	"game-catalog-backend/internal/httpapi"
	"game-catalog-backend/internal/media"
	sqlitestore "game-catalog-backend/internal/storage/sqlite"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type App struct {
	Server *echo.Echo
	store  *sqlitestore.Store
}

func newApp(dbPath string, seedDemoData bool) (*App, error) {
	return newAppWithBlobDir(dbPath, "blob", seedDemoData)
}

func newAppWithBlobDir(dbPath, blobDir string, seedDemoData bool) (*App, error) {
	store, err := sqlitestore.New(dbPath, seedDemoData)
	if err != nil {
		return nil, err
	}

	uploader, err := media.NewUploader(filepath.Clean(blobDir))
	if err != nil {
		_ = store.Close()
		return nil, err
	}

	gameService := game.NewService(store)
	server := httpapi.NewServer(gameService, uploader, filepath.Clean(blobDir))

	return &App{
		Server: server,
		store:  store,
	}, nil
}

func (a *App) Close() error {
	return a.store.Close()
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.Server.ServeHTTP(writer, request)
}
