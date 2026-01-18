package services

import (
	"context"
	"encoding/json"
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

	tables := []string{
		"ability_effects", "abilities", "unit_keywords", "weapons", "units", "factions",
		"keywords", "battle_formations", "enhancements", "rules", "games",
	}

	for _, table := range tables {
		_, err := dbpool.Exec(ctx, "DELETE FROM "+table)
		if err != nil {
			t.Logf("WARNING: failed to clear table %s: %v", table, err)
		}
	}

	cfg := &config.Config{Env: "test", Port: ":8080"}
	logger, _ := zap.NewDevelopment()

	return &state.State{
		DB:     queries,
		Cfg:    cfg,
		Logger: logger,
		Pool:   dbpool,
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
		GameID:             gameID,
		Name:               "Test Faction",
		Allegiance:         "Test Allegiance",
		Version:            "1.0",
		Source:             "Test Source",
		IsArmyOfRenown:     false,
		IsRegimentOfRenown: false,
		ParentFactionID:    uuid.NullUUID{},
	})
	if err != nil {
		t.Fatalf("failed to create faction: %v", err)
	}

	factionID := faction.ID

	return factionID
}

func createTestFactionWithName(t *testing.T, s *state.State, gameID uuid.UUID, name string) uuid.UUID {
	ctx := context.Background()

	faction, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:             gameID,
		Name:               name,
		Allegiance:         "Test Allegiance",
		Version:            "1.0",
		Source:             "Test Source",
		IsArmyOfRenown:     false,
		IsRegimentOfRenown: false,
		ParentFactionID:    uuid.NullUUID{},
	})
	if err != nil {
		t.Fatalf("failed to create faction with name: %v", err)
	}

	factionID := faction.ID
	return factionID
}

func createTestUnit(t *testing.T, s *state.State, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            "Test Unit",
		Description:     "",
		Move:            "10\"",
		HealthWounds:    "4",
		SaveStats:       "3+",
		WardFnp:         "6+",
		ControlOc:       "1",
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     4,
		MaxUnitSize:     8,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        false,
		AdditionalStats: json.RawMessage("{}"),
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create unit: %v", err)
	}
	unitID := unit.ID

	return unitID
}

func createTestUnitWithName(t *testing.T, s *state.State, factionID uuid.UUID, name string) uuid.UUID {
	ctx := context.Background()
	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            name,
		Description:     "",
		Move:            "10\"",
		HealthWounds:    "4",
		SaveStats:       "3+",
		WardFnp:         "6+",
		ControlOc:       "1",
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     4,
		MaxUnitSize:     8,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        false,
		AdditionalStats: json.RawMessage("{}"),
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create unit: %v", err)
	}

	unitID := unit.ID

	return unitID
}

func createTestUniqueUnit(t *testing.T, s *state.State, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            "Test Unique Unit",
		Description:     "",
		Move:            "10\"",
		HealthWounds:    "4",
		SaveStats:       "3+",
		WardFnp:         "6+",
		ControlOc:       "1",
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     4,
		MaxUnitSize:     8,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        true,
		AdditionalStats: json.RawMessage("{}"),
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create unique unit: %v", err)
	}
	unitID := unit.ID

	return unitID
}

func createTestWeapon(t *testing.T, s *state.State, unitID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	weapon, err := s.DB.CreateWeapon(ctx, database.CreateWeaponParams{
		UnitID:        unitID,
		Name:          "Test AoS Weapon",
		Range:         "10\"",
		Attacks:       "5",
		HitStats:      "3+",
		WoundStrength: "3+",
		RendAp:        "-1",
		Damage:        "2",
		Version:       "1.0",
		Source:        "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create weapon: %v", err)
	}

	weaponID := weapon.ID
	return weaponID
}

func createTestKeyword(t *testing.T, s *state.State, gameID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	keyword, err := s.DB.CreateKeyword(ctx, database.CreateKeywordParams{
		GameID:      gameID,
		Name:        "Test Keyword",
		Description: "Keyword Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create keyword: %v", err)
	}

	keywordID := keyword.ID

	return keywordID
}

func createTestAbilityUnit(t *testing.T, s *state.State, unitID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	ability, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID),
		FactionID: uuid.NullUUID{},
		Name:      "Test Ability",
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create unit ability: %v", err)
	}

	abilityID := ability.ID
	return abilityID
}

func createTestAbilityUnitWithName(t *testing.T, s *state.State, unitID uuid.UUID, name string) uuid.UUID {
	ctx := context.Background()

	ability, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID),
		FactionID: uuid.NullUUID{},
		Name:      name,
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create unit ability: %v", err)
	}

	abilityID := ability.ID
	return abilityID
}
