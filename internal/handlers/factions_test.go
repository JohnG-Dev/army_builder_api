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
	"github.com/google/uuid"
)

// setupTestFactionState clears factions table and inserts test data
func setupTestFactionState(t *testing.T) *state.State {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	// Clear tables (factions and games, since factions depend on a game_id)
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	// Insert test game
	var gameID uuid.UUID
	err = dbpool.QueryRow(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar','4e') RETURNING id;").Scan(&gameID)
	if err != nil {
		t.Fatalf("failed to insert game: %v", err)
	}

	// Insert test faction
	_, err = dbpool.Exec(ctx, "INSERT INTO factions (game_id, name) VALUES ($1,'Stormcast Eternals');", gameID)
	if err != nil {
		t.Fatalf("failed to insert faction: %v", err)
	}

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	return &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}
}

func TestListFactions_ReturnsFaction(t *testing.T) {
	s := setupTestFactionState(t)
	fHandlers := &FactionHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions", nil)
	w := httptest.NewRecorder()

	fHandlers.ListFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", res.StatusCode)
	}

	if !strings.Contains(string(body), "Stormcast Eternals") {
		t.Errorf("expected body to contain Stormcast, got %s", string(body))
	}
}

func TestGetFactionByID_NotFound(t *testing.T) {
	s := setupTestFactionState(t)
	fHandlers := &FactionHandlers{S: s}

	// Use a random UUID that won’t exist
	req := httptest.NewRequest(http.MethodGet, "/factions/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()

	// simulate param extraction (Go 1.22 mux replacement, for now just manually)
	req.SetPathValue("id", uuid.NewString()) // if using Go 1.22 path pattern

	fHandlers.GetFactionByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 got %d", res.StatusCode)
	}
}

func TestListFactions_Empty(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, _ := pgxpool.New(ctx, dbURL)
	queries := database.New(dbpool)

	// Truncate completely
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()
	s := &state.State{DB: queries, Cfg: cfg, Logger: logger}

	fHandlers := &FactionHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions", nil)
	w := httptest.NewRecorder()

	fHandlers.ListFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", res.StatusCode)
	}
	if string(body) != "[]" {
		t.Errorf("expected empty array, got %s", string(body))
	}
}

func TestListFactions_FilterByGameID(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}
	queries := database.New(dbpool)

	// Reset DB
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	// Insert two games
	var aosID, wh40kID uuid.UUID
	_ = dbpool.QueryRow(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar','4e') RETURNING id").Scan(&aosID)
	_ = dbpool.QueryRow(ctx, "INSERT INTO games (name, edition) VALUES ('Warhammer 40K','10e') RETURNING id").Scan(&wh40kID)

	// Insert factions for both games
	_, _ = dbpool.Exec(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Stormcast Eternals')", aosID)
	_, _ = dbpool.Exec(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Orruks')", aosID)
	_, _ = dbpool.Exec(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Space Marines')", wh40kID)

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()
	s := &state.State{DB: queries, Cfg: cfg, Logger: logger}
	fHandlers := &FactionHandlers{S: s}

	// ✅ Request with game_id filter (AoS)
	req := httptest.NewRequest(http.MethodGet, "/factions?game_id="+aosID.String(), nil)
	w := httptest.NewRecorder()
	fHandlers.ListFactions(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	// Factions from AoS should be present
	if !strings.Contains(bodyStr, "Stormcast Eternals") || !strings.Contains(bodyStr, "Orruks") {
		t.Errorf("expected AoS factions in body, got %s", bodyStr)
	}

	// 40K faction should NOT be present
	if strings.Contains(bodyStr, "Space Marines") {
		t.Errorf("expected no 40K factions, but got body %s", bodyStr)
	}
}
