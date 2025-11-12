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
			faction_id, name, points, move, health, save, ward, control,
			rend, attacks, damage, summon_cost, banishment,
			is_manifestation, min_size, max_size
		) VALUES ($1, 'Liberators', 100, '5"', 2, '4+', '6+', 2,
			'-', '2', '1', NULL, NULL, FALSE, 5, 15);`,
		factionID,
	)
	if err != nil {
		t.Fatalf("failed to insert unit: %v", err)
	}

	// Create test manifestation (same faction)
	_, err = dbpool.Exec(ctx, `
		INSERT INTO units (
			faction_id, name, points, move, health, save, ward, control,
			rend, attacks, damage, summon_cost, banishment,
			is_manifestation, min_size, max_size
		) VALUES ($1, 'Celestant Prime Manifestation', 200, '12"', 10, '3+', '5+', 0,
			'-2', '4', 'D6', '7+', '8+', TRUE, 1, 1);`,
		factionID,
	)
	if err != nil {
		t.Fatalf("failed to insert manifestation: %v", err)
	}

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	return &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
	}
}

// ✅ Normal unit list (non-manifestations)
func TestGetUnits_ReturnsUnits(t *testing.T) {
	s := setupTestUnitState(t)
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()
	uHandlers.GetNonManifestationUnits(w, req)

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
	if strings.Contains(body, "Celestant Prime Manifestation") {
		t.Errorf("expected not to include manifestations, got %s", body)
	}
}

// ✅ Manifestation list handler
func TestGetManifestations_ReturnsManifestations(t *testing.T) {
	s := setupTestUnitState(t)
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/manifestations", nil)
	w := httptest.NewRecorder()
	uHandlers.GetManifestations(w, req)

	res := w.Result()
	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", res.StatusCode)
	}
	if !strings.Contains(body, "Celestant Prime Manifestation") {
		t.Errorf("expected manifestation in response, got %s", body)
	}
	if strings.Contains(body, "Liberators") {
		t.Errorf("expected only manifestations, got %s", body)
	}
}

// ✅ Manifestation not found by ID
func TestGetManifestationByID_NotFound(t *testing.T) {
	s := setupTestUnitState(t)
	uHandlers := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/manifestations/"+uuid.NewString(), nil)
	req.SetPathValue("id", uuid.NewString())
	w := httptest.NewRecorder()

	uHandlers.GetManifestationByID(w, req)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 got %d", res.StatusCode)
	}
}

// ✅ Empty DB - ensure empty slices
func TestUnits_EmptyDB(t *testing.T) {
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

	checkEmpty := func(path string, fn func(http.ResponseWriter, *http.Request)) {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		fn(w, req)
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

	// Check both regular units and manifestations return empty slices
	checkEmpty("/units", uHandlers.GetNonManifestationUnits)
	checkEmpty("/manifestations", uHandlers.GetManifestations)
}

// ✅ Filter by faction_id
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

	var gameID uuid.UUID
	_ = dbpool.QueryRow(ctx, "INSERT INTO games (name, edition) VALUES ('Age of Sigmar','4e') RETURNING id").Scan(&gameID)

	var stormcastID, orrukID uuid.UUID
	_ = dbpool.QueryRow(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Stormcast Eternals') RETURNING id", gameID).Scan(&stormcastID)
	_ = dbpool.QueryRow(ctx, "INSERT INTO factions (game_id, name) VALUES ($1, 'Orruks') RETURNING id", gameID).Scan(&orrukID)

	// Units for both factions
	_, _ = dbpool.Exec(ctx, `
		INSERT INTO units (faction_id, name, points, move, health, save, ward, control, min_size, max_size, is_manifestation)
		VALUES ($1, 'Liberators', 100, '5"', 2, '4+', '6+', 2, 5, 15, FALSE)
	`, stormcastID)
	_, _ = dbpool.Exec(ctx, `
		INSERT INTO units (faction_id, name, points, move, health, save, ward, control, min_size, max_size, is_manifestation)
		VALUES ($1, 'Brutes', 150, '5"', 3, '4+', '6+', 3, 5, 10, FALSE)
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
	if !strings.Contains(body, "Liberators") {
		t.Errorf("expected Stormcast units, got %s", body)
	}
	if strings.Contains(body, "Brutes") {
		t.Errorf("expected no Orruk units, got %s", body)
	}
}
