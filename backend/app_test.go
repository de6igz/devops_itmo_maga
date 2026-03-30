package main

import (
	"bytes"
	"encoding/json"
	"game-catalog-backend/internal/game"
	"game-catalog-backend/internal/httpapi"
	"game-catalog-backend/internal/media"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
)

func setupTestServer(t *testing.T, repository game.Repository) http.Handler {
	t.Helper()

	blobDir := filepath.Join(t.TempDir(), "blob")
	uploader, err := media.NewUploader(blobDir)
	if err != nil {
		t.Fatalf("failed to create uploader: %v", err)
	}

	service := game.NewService(repository)
	return httpapi.NewServer(service, uploader, blobDir)
}

func performMultipartRequest(t *testing.T, handler http.Handler, method, path string, fieldName, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create multipart file: %v", err)
	}

	if _, err := io.Copy(fileWriter, bytes.NewReader(content)); err != nil {
		t.Fatalf("failed to copy multipart file: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func performRequest(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
	}

	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateGame(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := NewMockGameRepository(controller)
	server := setupTestServer(t, repository)

	repository.EXPECT().
		Create(gomock.AssignableToTypeOf(game.Game{})).
		DoAndReturn(func(entity game.Game) (game.Game, error) {
			if entity.Title != "Cyberpunk 2077" {
				t.Fatalf("unexpected title: %s", entity.Title)
			}
			if entity.Description == "" {
				t.Fatalf("description should not be empty")
			}

			entity.ID = 1
			return entity, nil
		})

	response := performRequest(t, server, http.MethodPost, "/api/games", map[string]any{
		"title":       "Cyberpunk 2077",
		"genre":       "RPG",
		"platform":    "PC",
		"releaseYear": 2020,
		"rating":      8,
		"status":      "playing",
		"description": "Футуристическая RPG с открытым городом, сюжетом и вариативной прокачкой.",
	})

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}
}

func TestGetGames(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := NewMockGameRepository(controller)
	server := setupTestServer(t, repository)

	repository.EXPECT().
		List(game.Filters{}).
		Return([]game.Game{
			{
				ID:          1,
				Title:       "Celeste",
				Genre:       "Platformer",
				Platform:    "Switch",
				ReleaseYear: 2018,
				Rating:      9,
				Status:      "completed",
				Description: "Точная платформенная игра про преодоление и внутренний рост.",
			},
		}, nil)

	response := performRequest(t, server, http.MethodGet, "/api/games", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var games []game.Game
	if err := json.NewDecoder(response.Body).Decode(&games); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(games) != 1 || games[0].Title != "Celeste" {
		t.Fatalf("unexpected games list: %+v", games)
	}
}

func TestUpdateGame(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := NewMockGameRepository(controller)
	server := setupTestServer(t, repository)

	repository.EXPECT().
		Update(int64(1), gomock.AssignableToTypeOf(game.Game{})).
		DoAndReturn(func(id int64, entity game.Game) (game.Game, error) {
			if id != 1 {
				t.Fatalf("unexpected id: %d", id)
			}
			if entity.Title != "Control Ultimate Edition" {
				t.Fatalf("unexpected title: %s", entity.Title)
			}

			entity.ID = id
			return entity, nil
		})

	response := performRequest(t, server, http.MethodPut, "/api/games/1", map[string]any{
		"title":       "Control Ultimate Edition",
		"genre":       "Action",
		"platform":    "PC",
		"releaseYear": 2021,
		"rating":      9,
		"status":      "completed",
		"description": "Полное издание атмосферного action с дополнениями и необычной боевой системой.",
	})

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestDeleteGame(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := NewMockGameRepository(controller)
	server := setupTestServer(t, repository)

	repository.EXPECT().
		Delete(int64(1)).
		Return(nil)

	response := performRequest(t, server, http.MethodDelete, "/api/games/1", nil)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
}

func TestInvalidGameData(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := NewMockGameRepository(controller)
	server := setupTestServer(t, repository)

	response := performRequest(t, server, http.MethodPost, "/api/games", map[string]any{
		"title":       "",
		"genre":       "RPG",
		"platform":    "PC",
		"releaseYear": 1900,
		"rating":      20,
		"status":      "unknown",
		"description": "",
	})

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestUploadImage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := NewMockGameRepository(controller)
	server := setupTestServer(t, repository)

	response := performMultipartRequest(
		t,
		server,
		http.MethodPost,
		"/api/uploads/image",
		"image",
		"cover.png",
		[]byte("fake image content"),
	)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode upload response: %v", err)
	}

	if payload["path"] == "" {
		t.Fatalf("expected image path in response")
	}
}
