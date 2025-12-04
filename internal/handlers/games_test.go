package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

// setupTestDB creates a fresh DB connection and clears the games table
func setupTestDB(t *testing.T) *state.State {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	// Clear all tables to ensure clean state
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	return &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}
}

func TestGetGames_ReturnsGame(t *testing.T) {
	s := setupTestDB(t)

	// Insert test data
	ctx := context.Background()
	_, err := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Age of Sigmar",
		Edition: "4e",
	})
	if err != nil {
		t.Fatalf("failed to insert game: %v", err)
	}

	// Create handler
	handler := &GamesHandlers{S: s}

	// Make request
	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	handler.GetGames(w, req)

	// Assert
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Age of Sigmar") {
		t.Errorf("expected body to contain 'Age of Sigmar', got %s", bodyStr)
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
	if string(body) != "[]" {
		t.Errorf("expected empty array [], got %s", string(body))
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
		t.Errorf("expected both games in response, got %s", bodyStr)
	}
}

func TestGetGameByName(t *testing.T) {
	s := setupTestDB(t)

	ctx := context.Background()
	_, _ = s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Age of Sigmar",
		Edition: "4e",
	})

	handler := &GamesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/games?name=Age of Sigmar", nil)
	w := httptest.NewRecorder()

	handler.GetGames(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Age of Sigmar") {
		t.Errorf("expected to find game by name, got %s", string(body))
	}
}
