package httpapi

import (
	"net/http"
	"strconv"

	"game-catalog-backend/internal/game"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  GameService
	uploader ImageUploader
}

func (h *Handler) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listGames(c echo.Context) error {
	games, err := h.service.List(game.Filters{
		Status: c.QueryParam("status"),
		Genre:  c.QueryParam("genre"),
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, games)
}

func (h *Handler) getGame(c echo.Context) error {
	id, err := parseGameID(c.Param("id"))
	if err != nil {
		return err
	}

	entity, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entity)
}

func (h *Handler) createGame(c echo.Context) error {
	var payload game.Game
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "некорректный JSON")
	}

	createdGame, err := h.service.Create(payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdGame)
}

func (h *Handler) updateGame(c echo.Context) error {
	id, err := parseGameID(c.Param("id"))
	if err != nil {
		return err
	}

	var payload game.Game
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "некорректный JSON")
	}

	updatedGame, err := h.service.Update(id, payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedGame)
}

func (h *Handler) deleteGame(c echo.Context) error {
	id, err := parseGameID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.Delete(id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) uploadImage(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "изображение не передано")
	}

	imagePath, err := h.uploader.SaveImage(fileHeader)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"path": imagePath})
}

func parseGameID(rawID string) (int64, error) {
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "некорректный идентификатор")
	}

	return id, nil
}
