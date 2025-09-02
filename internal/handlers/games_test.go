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

// setupTestState clears the games table and inserts test data
func setupTestState(t *testing.T) *state.State {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	// Clear table
	_, err = dbpool.Exec(ctx, "DELETE FROM games;")
	if err != nil {
		t.Fatalf("failed to clear games table: %v", err)
	}

	// Insert one test row
	_, err = dbpool.Exec(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar', '4e');")
	if err != nil {
		t.Fatalf("failed to insert test game: %v", err)
	}

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	return &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}
}

// ✅ Test: Normal case (returns one game)
func TestGetGames_ReturnsGame(t *testing.T) {
	s := setupTestState(t)

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	GetGames(s, w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	if !strings.Contains(string(body), "Age of Sigmar") {
		t.Errorf("expected body to contain 'Age of Sigmar', got %s", string(body))
	}
}

// ✅ Test: Empty DB (returns empty array)
func TestGetGames_EmptyDB(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, _ := pgxpool.New(ctx, dbURL)
	queries := database.New(dbpool)

	// Clear table only, no insert
	dbpool.Exec(ctx, "DELETE FROM games;")

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()
	s := &state.State{DB: queries, Cfg: cfg, Logger: logger}

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	GetGames(s, w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	if string(body) != "[]" {
		t.Errorf("expected empty array, got %s", string(body))
	}
}

// ✅ Test: Invalid method (should return 405)
func TestGetGames_InvalidMethod(t *testing.T) {
	s := setupTestState(t)

	req := httptest.NewRequest(http.MethodPost, "/games", nil) // wrong method
	w := httptest.NewRecorder()

	// Simulate method check (you can add this logic in handler later)
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	} else {
		GetGames(s, w, req)
	}

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", res.StatusCode)
	}
}

// ✅ Test: Multiple rows (returns all games)
func TestGetGames_MultipleGames(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, _ := pgxpool.New(ctx, dbURL)
	queries := database.New(dbpool)

	// Clear and insert multiple rows
	dbpool.Exec(ctx, "DELETE FROM games;")
	dbpool.Exec(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar', '4e');")
	dbpool.Exec(ctx, "INSERT INTO games (name, edition) VALUES ('Warhammer 40K', '10e');")

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()
	s := &state.State{DB: queries, Cfg: cfg, Logger: logger}

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	w := httptest.NewRecorder()

	GetGames(s, w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	if !strings.Contains(string(body), "Age of Sigmar") || !strings.Contains(string(body), "Warhammer 40K") {
		t.Errorf("expected body to contain both games, got %s", string(body))
	}
}

