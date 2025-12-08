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
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/config"
	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

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

func createTestGame(t *testing.T, s *state.State) uuid.UUID {
	ctx := context.Background()

	game, err := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Test Game",
		Edition: "Test Edition",
		Version: "1.0",
		Source:  "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	gameID := game.ID

	return gameID
}

func createTestFaction(t *testing.T, s *state.State, gameID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	faction, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:  gameID,
		Name:    "Test Faction",
		Version: "1.0",
		Source:  "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create faction: %v", err)
	}

	factionID := faction.ID

	return factionID
}

func createTestUnit(t *testing.T, s *state.State, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            "Test Unit",
		Move:            10,
		Health:          4,
		Save:            "3+",
		Ward:            "6+",
		Control:         1,
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     4,
		MaxUnitSize:     8,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        false,
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create unit: %v", err)
	}
	unitID := unit.ID

	return unitID
}

func createTestWeapon(t *testing.T, s *state.State, unitID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	weapon, err := s.DB.CreateWeapon(ctx, database.CreateWeaponParams{
		UnitID:  unitID,
		Name:    "Test AoS Weapon",
		Range:   "10\"",
		Attacks: "5",
		ToHit:   "3+",
		ToWound: "3+",
		Rend:    "-1",
		Damage:  "2",
		Version: "1.0",
		Source:  "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create weapon: %v", err)
	}

	weaponID := weapon.ID
	return weaponID
}
