package handlers

import (
	"context"
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

func createTestGameWithName(t *testing.T, s *state.State, name string) uuid.UUID {
	ctx := context.Background()
	if name == "" {
		t.Fatalf("game name required")
	}

	game, err := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    name,
		Edition: "Test Edition",
		Version: "1.0",
		Source:  "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create game with name: %v", err)
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

func createTestFactionWithName(t *testing.T, s *state.State, gameID uuid.UUID, name string) uuid.UUID {
	ctx := context.Background()

	faction, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:  gameID,
		Name:    name,
		Version: "1.0",
		Source:  "Test Source",
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

func createTestUnitWithName(t *testing.T, s *state.State, factionID uuid.UUID, name string) uuid.UUID {
	ctx := context.Background()
	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            name,
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

func createTestManifestation(t *testing.T, s *state.State, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       factionID,
		Name:            "Test Manifestation",
		Move:            10,
		Health:          7,
		Save:            "5+",
		Ward:            "6+",
		Control:         1,
		Points:          100,
		SummonCost:      "7",
		Banishment:      "5",
		MinUnitSize:     4,
		MaxUnitSize:     8,
		MatchedPlay:     true,
		IsManifestation: true,
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

func createTestRule(t *testing.T, s *state.State, gameID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	rule, err := s.DB.CreateRule(ctx, database.CreateRuleParams{
		GameID:      gameID,
		Name:        "Test Rule",
		RuleType:    "core",
		Text:        "Rule text content",
		Description: "Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	ruleID := rule.ID

	return ruleID
}

func createTestBattleFormation(t *testing.T, s *state.State, gameID uuid.UUID, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	battleFormation, err := s.DB.CreateBattleFormation(ctx, database.CreateBattleFormationParams{
		GameID:      gameID,
		FactionID:   factionID,
		Name:        "Test BattleFormation",
		Description: "Formation Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create battle formation: %v", err)
	}

	battleFormationID := battleFormation.ID
	return battleFormationID
}

func createTestEnhancement(t *testing.T, s *state.State, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	enhancement, err := s.DB.CreateEnhancement(ctx, database.CreateEnhancementParams{
		FactionID:       factionID,
		Name:            "Test Enhancement",
		EnhancementType: "artefact",
		Description:     "Enhancement Description",
		Points:          20,
		IsUnique:        true,
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create enhancement: %v", err)
	}

	enhancementID := enhancement.ID

	return enhancementID
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

func createTestKeywordWithName(t *testing.T, s *state.State, gameID uuid.UUID, name string) uuid.UUID {
	ctx := context.Background()

	keyword, err := s.DB.CreateKeyword(ctx, database.CreateKeywordParams{
		GameID:      gameID,
		Name:        name,
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

func createTestAbilityFaction(t *testing.T, s *state.State, factionID uuid.UUID) uuid.UUID {
	ctx := context.Background()

	ability, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    uuid.NullUUID{},
		FactionID: database.UUIDToNullUUID(factionID),
		Name:      "Test Ability",
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create faction ability: %v", err)
	}

	abilityID := ability.ID
	return abilityID
}
