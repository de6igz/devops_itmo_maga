package httpapi

import (
	"errors"
	"net/http"

	"game-catalog-backend/internal/game"

	"github.com/labstack/echo/v4"
)

func mapError(err error) (int, string) {
	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		if message, ok := httpError.Message.(string); ok {
			return httpError.Code, message
		}
		return httpError.Code, http.StatusText(httpError.Code)
	}

	var validationError game.ValidationError
	if errors.As(err, &validationError) {
		return http.StatusBadRequest, validationError.Error()
	}

	if errors.Is(err, game.ErrNotFound) {
		return http.StatusNotFound, "игра не найдена"
	}

	return http.StatusInternalServerError, "внутренняя ошибка сервера"
}
