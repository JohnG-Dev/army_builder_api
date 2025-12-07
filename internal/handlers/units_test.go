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

// setupTestUnitDB creates fresh DB with game + faction + units
func setupTestUnitDB(t *testing.T) (*state.State, uuid.UUID) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	queries := database.New(dbpool)

	// Clear tables in dependency order
	_, _ = dbpool.Exec(ctx, "DELETE FROM units;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	// Create game
	game, _ := queries.CreateGame(ctx, database.CreateGameParams{
		Name:    "Age of Sigmar",
		Edition: "4e",
	})

	// Create faction
	faction, _ := queries.CreateFaction(ctx, database.CreateFactionParams{
		GameID: game.ID,
		Name:   "Stormcast Eternals",
	})

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	state := &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}

	return state, faction.ID
}

func TestGetUnits_ReturnsUnits(t *testing.T) {
	s, factionID := setupTestUnitDB(t)

	ctx := context.Background()

	// Insert unit
	_, _ = s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:   factionID,
		Name:        "Liberators",
		MinUnitSize: 5,
		MaxUnitSize: 15,
		Points:      95,
		Move:        5,
		Health:      2,
		Save:        "3+",
		Ward:        "—",
		Control:     2,
		MatchedPlay: true,
	})

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()

	handler.GetNonManifestationUnits(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Liberators") {
		t.Errorf("expected unit in response, got %s", string(body))
	}
}

func TestGetUnits_EmptyDB(t *testing.T) {
	s, _ := setupTestUnitDB(t)

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()

	handler.GetNonManifestationUnits(w, req)

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

func TestGetUnits_FilterByFactionID(t *testing.T) {
	s, factionID := setupTestUnitDB(t)

	ctx := context.Background()

	// Create second faction
	game, _ := s.DB.GetGames(ctx)
	faction2, _ := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID: game[0].ID,
		Name:   "Orruks",
	})

	// Insert units for both factions
	_, _ = s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:   factionID,
		Name:        "Liberators",
		Points:      95,
		MinUnitSize: 5,
		MaxUnitSize: 15,
	})
	_, _ = s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:   faction2.ID,
		Name:        "Brutes",
		Points:      150,
		MinUnitSize: 5,
		MaxUnitSize: 10,
	})

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units?faction_id="+factionID.String(), nil)
	w := httptest.NewRecorder()

	handler.GetUnits(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Liberators") {
		t.Errorf("expected Stormcast unit, got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Brutes") {
		t.Errorf("expected no Orruk units, got %s", bodyStr)
	}
}

func TestGetManifestations_ReturnsManifestations(t *testing.T) {
	s, factionID := setupTestUnitDB(t)

	ctx := context.Background()

	// Insert manifestation
	_, _ = s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            "Celestant Prime Manifestation",
		IsManifestation: true,
		Points:          200,
		MinUnitSize:     1,
		MaxUnitSize:     1,
	})

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/manifestations", nil)
	w := httptest.NewRecorder()

	handler.GetManifestations(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Celestant Prime Manifestation") {
		t.Errorf("expected manifestation in response, got %s", string(body))
	}
}

func TestGetUnitByID_Success(t *testing.T) {
	s, factionID := setupTestUnitDB(t)

	ctx := context.Background()

	unit, _ := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:   factionID,
		Name:        "Liberators",
		Points:      95,
		MinUnitSize: 5,
		MaxUnitSize: 15,
	})

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units/"+unit.ID.String(), nil)
	req.SetPathValue("id", unit.ID.String())
	w := httptest.NewRecorder()

	handler.GetUnitByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Liberators") {
		t.Errorf("expected unit in response, got %s", string(body))
	}
}

func TestGetUnitByID_NotFound(t *testing.T) {
	s, _ := setupTestUnitDB(t)

	handler := &UnitsHandlers{S: s}

	randomID := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/units/"+randomID, nil)
	req.SetPathValue("id", randomID)
	w := httptest.NewRecorder()

	handler.GetUnitByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", res.StatusCode)
	}
}
