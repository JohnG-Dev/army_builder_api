package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

// setupTestFactionDB creates fresh DB with game + faction
func setupTestFactionDB(t *testing.T) (*state.State, uuid.UUID) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	// Clear tables
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	// Insert game
	game, err := queries.CreateGame(ctx, database.CreateGameParams{
		Name:    "Age of Sigmar",
		Edition: "4e",
	})
	if err != nil {
		t.Fatalf("failed to insert game: %v", err)
	}

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	state := &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}

	return state, game.ID
}

func TestListFactions_ReturnsFaction(t *testing.T) {
	s, gameID := setupTestFactionDB(t)

	ctx := context.Background()

	// Insert faction
	_, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:     gameID,
		Name:       "Stormcast Eternals",
		Allegiance: "ORDER",
	})
	if err != nil {
		t.Fatalf("failed to insert faction: %v", err)
	}

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions", nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Stormcast Eternals") {
		t.Errorf("expected faction in response, got %s", string(body))
	}
}

func TestListFactions_Empty(t *testing.T) {
	s, _ := setupTestFactionDB(t)

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions", nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)

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

func TestListFactions_FilterByGameID(t *testing.T) {
	s, gameID := setupTestFactionDB(t)

	ctx := context.Background()

	// Create second game
	game2, _ := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Warhammer 40K",
		Edition: "10e",
	})

	// Insert factions for both games
	_, _ = s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID: gameID,
		Name:   "Stormcast Eternals",
	})
	_, _ = s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID: game2.ID,
		Name:   "Space Marines",
	})

	handler := &FactionsHandlers{S: s}

	// Request with game_id filter
	req := httptest.NewRequest(http.MethodGet, "/factions?game_id="+gameID.String(), nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Stormcast Eternals") {
		t.Errorf("expected AoS faction, got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Space Marines") {
		t.Errorf("expected no 40K factions, got %s", bodyStr)
	}
}

func TestGetFactionByID_NotFound(t *testing.T) {
	s, _ := setupTestFactionDB(t)

	handler := &FactionsHandlers{S: s}

	randomID := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/factions/"+randomID, nil)
	req.SetPathValue("id", randomID)
	w := httptest.NewRecorder()

	handler.GetFactionByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", res.StatusCode)
	}
}

func TestGetFactionByID_Success(t *testing.T) {
	s, gameID := setupTestFactionDB(t)

	ctx := context.Background()

	faction, _ := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID: gameID,
		Name:   "Stormcast Eternals",
	})

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions/"+faction.ID.String(), nil)
	req.SetPathValue("id", faction.ID.String())
	w := httptest.NewRecorder()

	handler.GetFactionByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Stormcast Eternals") {
		t.Errorf("expected faction in response, got %s", string(body))
	}
}
