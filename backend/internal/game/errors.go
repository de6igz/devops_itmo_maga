package game

import "errors"

var ErrNotFound = errors.New("game not found")

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
