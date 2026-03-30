package httpapi

import (
	"mime/multipart"
	"net/http"

	"game-catalog-backend/internal/game"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GameService interface {
	List(game.Filters) ([]game.Game, error)
	GetByID(int64) (game.Game, error)
	Create(game.Game) (game.Game, error)
	Update(int64, game.Game) (game.Game, error)
	Delete(int64) error
}

type ImageUploader interface {
	SaveImage(fileHeader *multipart.FileHeader) (string, error)
}

func NewServer(service GameService, uploader ImageUploader, blobDir string) *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Use(middleware.Recover())

	handler := &Handler{service: service, uploader: uploader}
	api := server.Group("/api")

	api.GET("/health", handler.health)
	api.GET("/games", handler.listGames)
	api.GET("/games/:id", handler.getGame)
	api.POST("/games", handler.createGame)
	api.PUT("/games/:id", handler.updateGame)
	api.DELETE("/games/:id", handler.deleteGame)
	api.POST("/uploads/image", handler.uploadImage)

	server.Static("/blob", blobDir)

	server.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		status, message := mapError(err)
		_ = c.JSON(status, map[string]string{"error": message})
	}

	server.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return server
}
