package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

func TestGetGames_ReturnsGame(t *testing.T) {
	s := setupTestDB(t)

	createTestGame(t, s)

	handler := &GamesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	handler.GetGames(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Game") {
		t.Errorf("expected body to contain 'Test Game', got %s", bodyStr)
	}
}

func TestGetGames_EmptyDB(t *testing.T) {
	s := setupTestDB(t)

	handler := &GamesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	handler.GetGames(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)
	if bodyStr != "[]" {
		t.Errorf("expected empty array [], got %s", bodyStr)
	}
}

func TestGetGames_MultipleGames(t *testing.T) {
	s := setupTestDB(t)

	ctx := context.Background()

	// Insert multiple games
	_, _ = s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Age of Sigmar",
		Edition: "4e",
	})
	_, _ = s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Warhammer 40K",
		Edition: "10e",
	})

	handler := &GamesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	handler.GetGames(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Age of Sigmar") || !strings.Contains(bodyStr, "Warhammer 40K") {
		t.Errorf("expected body to contain 'Age of Sigmar' 'Warhammer 40k', got %s", bodyStr)
	}
}

func TestGetGameByName(t *testing.T) {
	s := setupTestDB(t)

	createTestGame(t, s)

	handler := &GamesHandlers{S: s}

	gameName := "Test Game"
	encodedName := url.QueryEscape(gameName)

	req := httptest.NewRequest(http.MethodGet, "/games?name="+encodedName, nil)
	w := httptest.NewRecorder()

	handler.GetGames(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "Test Game") {
		t.Errorf("expected body to contain 'Test Game', got %s", bodyStr)
	}
}

func TestGamesByID_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)

	handler := &GamesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/games/"+gameID.String(), nil)
	req.SetPathValue("id", gameID.String())

	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Game") {
		t.Errorf("expected body to contain 'Test Game', got %s", bodyStr)
	}
}

func TestGameByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	handler := &GamesHandlers{S: s}

	nonExistentID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/games/"+nonExistentID.String(), nil)
	req.SetPathValue("id", nonExistentID.String())
	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}
}

func TestGameByID_InvalidUUID(t *testing.T) {
	s := setupTestDB(t)

	handler := &GamesHandlers{S: s}

	invalidID := "not-a-valid-uuid"

	req := httptest.NewRequest(http.MethodGet, "/games/"+invalidID, nil)
	req.SetPathValue("id", invalidID)
	w := httptest.NewRecorder()

	handler.GetGameByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", res.StatusCode)
	}
}
