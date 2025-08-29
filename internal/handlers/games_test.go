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

func setupTestState(t *testing.T) *state.State {

	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	_, err = dbpool.Exec(ctx, "DELETE FROM games;")
	if err != nil {
		t.Fatalf("failed to clear games table: %v", err)
	}

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

func TestGetGames(t *testing.T) {

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
