package game

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

var validStatuses = map[string]struct{}{
	"planned":   {},
	"playing":   {},
	"completed": {},
}

func Normalize(game *Game) {
	game.Title = strings.TrimSpace(game.Title)
	game.Genre = strings.TrimSpace(game.Genre)
	game.Platform = strings.TrimSpace(game.Platform)
	game.Status = strings.TrimSpace(game.Status)
}

func Validate(game Game) error {
	currentYear := time.Now().Year() + 1

	if game.Title == "" {
		return ValidationError{Message: "поле title обязательно"}
	}

	if game.Genre == "" {
		return ValidationError{Message: "поле genre обязательно"}
	}

	if game.Platform == "" {
		return ValidationError{Message: "поле platform обязательно"}
	}

	if game.ReleaseYear < 1970 || game.ReleaseYear > currentYear {
		return ValidationError{Message: fmt.Sprintf("поле releaseYear должно быть числом от 1970 до %d", currentYear)}
	}

	if game.Rating < 1 || game.Rating > 10 {
		return ValidationError{Message: "поле rating должно быть числом от 1 до 10"}
	}

	if _, ok := validStatuses[game.Status]; !ok {
		return ValidationError{Message: "поле status должно быть одним из: planned, playing, completed"}
	}

	if game.ImagePath != "" && !strings.HasPrefix(filepath.ToSlash(game.ImagePath), "/blob/") {
		return ValidationError{Message: "поле imagePath должно содержать путь внутри /blob/"}
	}

	return nil
}
