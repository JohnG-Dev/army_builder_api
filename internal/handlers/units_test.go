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

// setupTestUnitState clears relevant tables and inserts base data for units
func setupTestUnitState(t *testing.T) *state.State {
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
	var gameID uuid.UUID
	err = dbpool.QueryRow(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar','4e') RETURNING id;").Scan(&gameID)
	if err != nil {
		t.Fatalf("failed to insert game: %v", err)
	}

	// Create faction
	var factionID uuid.UUID
	err = dbpool.QueryRow(ctx, "INSERT INTO factions (game_id, name) VALUES ($1,'Stormcast Eternals') RETURNING id;", gameID).Scan(&factionID)
	if err != nil {
		t.Fatalf("failed to insert faction: %v", err)
	}

	// Create test unit
	_, err = dbpool.Exec(ctx, `
		INSERT INTO units (
			faction_id, name, points, move, health, save, ward, control, min_size, max_size
		) VALUES ($1, 'Liberators', 100, 5, 2, 4, 6, 2, 5, 15);
	`, factionID)
	if err != nil {
		t.Fatalf("failed to insert unit: %v", err)
	}

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	return &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}
}

func TestGetUnits_ReturnsUnits(t *testing.T) {
	s := setupTestUnitState(t)
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()
	uHandlers.GetUnits(w, req)

	res := w.Result()
	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", res.StatusCode)
	}
	if !strings.Contains(body, "Liberators") {
		t.Errorf("expected Liberators in response, got %s", body)
	}
}

func TestGetUnitByID_NotFound(t *testing.T) {
	s := setupTestUnitState(t)
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units/"+uuid.NewString(), nil)
	req.SetPathValue("id", uuid.NewString())
	w := httptest.NewRecorder()

	uHandlers.GetUnitByID(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 got %d", res.StatusCode)
	}
}

func TestGetUnits_Empty(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, _ := pgxpool.New(ctx, dbURL)
	queries := database.New(dbpool)

	_, _ = dbpool.Exec(ctx, "DELETE FROM units;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()
	s := &state.State{DB: queries, Cfg: cfg, Logger: logger}
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()
	uHandlers.GetUnits(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", res.StatusCode)
	}
	if string(body) != "[]" {
		t.Errorf("expected [], got %s", string(body))
	}
}

func TestGetUnits_FilterByFactionID(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable"
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}
	queries := database.New(dbpool)

	// Reset DB
	_, _ = dbpool.Exec(ctx, "DELETE FROM units;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM factions;")
	_, _ = dbpool.Exec(ctx, "DELETE FROM games;")

	var aosID uuid.UUID
	_ = dbpool.QueryRow(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar','4e') RETURNING id").Scan(&aosID)

	var stormcastID, orrukID uuid.UUID
	_ = dbpool.QueryRow(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Stormcast Eternals') RETURNING id", aosID).Scan(&stormcastID)
	_ = dbpool.QueryRow(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Orruks') RETURNING id", aosID).Scan(&orrukID)

	// Units for both factions
	_, _ = dbpool.Exec(ctx, `
		INSERT INTO units (faction_id, name, points, move, health, save, ward, control, min_size, max_size)
		VALUES ($1, 'Liberators', 100, 5, 2, 4, 6, 2, 5, 15);
	`, stormcastID)
	_, _ = dbpool.Exec(ctx, `
		INSERT INTO units (faction_id, name, points, move, health, save, ward, control, min_size, max_size)
		VALUES ($1, 'Judicators', 120, 5, 2, 4, 6, 2, 5, 15);
	`, stormcastID)
	_, _ = dbpool.Exec(ctx, `
		INSERT INTO units (faction_id, name, points, move, health, save, ward, control, min_size, max_size)
		VALUES ($1, 'Brutes', 150, 5, 3, 4, 6, 3, 5, 10);
	`, orrukID)

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()
	s := &state.State{DB: queries, Cfg: cfg, Logger: logger}
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units?faction_id="+stormcastID.String(), nil)
	w := httptest.NewRecorder()
	uHandlers.GetUnits(w, req)

	res := w.Result()
	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}
	if !strings.Contains(body, "Liberators") || !strings.Contains(body, "Judicators") {
		t.Errorf("expected Stormcast units, got %s", body)
	}
	if strings.Contains(body, "Brutes") {
		t.Errorf("expected no Orruk units, got %s", body)
	}
}

