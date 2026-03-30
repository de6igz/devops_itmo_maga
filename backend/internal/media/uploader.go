package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"game-catalog-backend/internal/game"
)

type Uploader struct {
	rootDir string
}

func NewUploader(rootDir string) (*Uploader, error) {
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return nil, err
	}

	return &Uploader{rootDir: rootDir}, nil
}

func (u *Uploader) SaveImage(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", game.ValidationError{Message: "файл изображения обязателен"}
	}

	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch extension {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return "", game.ValidationError{Message: "допустимы только изображения jpg, jpeg, png или webp"}
	}

	source, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	targetPath := filepath.Join(u.rootDir, fileName)

	target, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return "", err
	}

	return "/blob/" + fileName, nil
}
