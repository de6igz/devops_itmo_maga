package main

import (
	"bytes"
	"encoding/json"
	"game-catalog-backend/internal/game"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func setupTestApp(t *testing.T) (*App, func()) {
	t.Helper()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	blobDir := filepath.Join(tempDir, "blob")

	app, err := newAppWithBlobDir(dbPath, blobDir, false)
	if err != nil {
		t.Fatalf("failed to create test app: %v", err)
	}

	return app, func() {
		_ = app.Close()
	}
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
	app, cleanup := setupTestApp(t)
	defer cleanup()

	response := performRequest(t, app, http.MethodPost, "/api/games", map[string]any{
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
	app, cleanup := setupTestApp(t)
	defer cleanup()

	performRequest(t, app, http.MethodPost, "/api/games", map[string]any{
		"title":       "Celeste",
		"genre":       "Platformer",
		"platform":    "Switch",
		"releaseYear": 2018,
		"rating":      9,
		"status":      "completed",
		"description": "Точная платформенная игра про преодоление и внутренний рост.",
	})

	response := performRequest(t, app, http.MethodGet, "/api/games", nil)
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
	app, cleanup := setupTestApp(t)
	defer cleanup()

	performRequest(t, app, http.MethodPost, "/api/games", map[string]any{
		"title":       "Control",
		"genre":       "Action",
		"platform":    "PC",
		"releaseYear": 2019,
		"rating":      8,
		"status":      "planned",
		"description": "Мистический action с аномалиями, перестрелками и исследованием Бюро.",
	})

	response := performRequest(t, app, http.MethodPut, "/api/games/1", map[string]any{
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
	app, cleanup := setupTestApp(t)
	defer cleanup()

	performRequest(t, app, http.MethodPost, "/api/games", map[string]any{
		"title":       "Portal 2",
		"genre":       "Puzzle",
		"platform":    "PC",
		"releaseYear": 2011,
		"rating":      10,
		"status":      "completed",
		"description": "Головоломка от первого лица с порталами, кооперативом и ярким юмором.",
	})

	response := performRequest(t, app, http.MethodDelete, "/api/games/1", nil)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
}

func TestInvalidGameData(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	response := performRequest(t, app, http.MethodPost, "/api/games", map[string]any{
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
	app, cleanup := setupTestApp(t)
	defer cleanup()

	response := performMultipartRequest(
		t,
		app,
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
