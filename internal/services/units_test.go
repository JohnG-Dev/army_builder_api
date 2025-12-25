package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
)

func TestGetUnitByID_DeepMapping(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID1 := createTestUnitWithName(t, s, factionID, "Judicator")
	unitID2 := createTestUnitWithName(t, s, factionID, "Liberator")
	weaponID1 := createTestWeapon(t, s, unitID1)
	weaponID2 := createTestWeapon(t, s, unitID1)
	abilityID := createTestAbilityUnit(t, s, unitID1)
	keywordID := createTestKeyword(t, s, gameID)

	err := s.DB.AddKeywordToUnit(ctx, database.AddKeywordToUnitParams{
		UnitID:    unitID1,
		KeywordID: keywordID,
		Value:     "Wizard",
	})
	if err != nil {
		t.Fatalf("failed to link keyword to unit: %v", err)
	}

	unit1, err := GetUnitByID(s, ctx, unitID1)
	if err != nil {
		t.Fatalf("failed to get unit 1, %v", err)
	}
	unit2, err := GetUnitByID(s, ctx, unitID2)
	if err != nil {
		t.Fatalf("failed to get unit 2, %v", err)
	}

	if len(unit1.Weapons) != 2 {
		t.Errorf("service failed to create weapons: expected 2, got %d", len(unit1.Weapons))
	}

	if unit1.Weapons[0].ID != weaponID1 {
		t.Errorf("expected weapon id %v, got %v", weaponID1, unit1.Weapons[0].ID)
	}

	if unit1.Weapons[1].ID != weaponID2 {
		t.Errorf("expected weapon id %v, got %v", weaponID2, unit1.Weapons[1].ID)
	}

	if len(unit1.Keywords) != 1 {
		t.Errorf("service failed to create keyword: expected 1, got %d", len(unit1.Keywords))
	}

	if len(unit1.Abilities) != 1 {
		t.Errorf("service failed to create ability, expected 1, got %d", len(unit1.Abilities))
	}

	if unit1.Abilities[0].ID != abilityID {
		t.Errorf("expected ability id %v, got %v", abilityID, unit1.Abilities[0].ID)
	}

	if unit1.Name != "Judicator" {
		t.Errorf("expected name: Judicator, got %s", unit1.Name)
	}

	if len(unit2.Weapons) != 0 {
		t.Errorf("servce linked weapon to incorrect unit: expected 0, got %d", len(unit2.Weapons))
	}

	if len(unit2.Abilities) != 0 {
		t.Errorf("service linked ability to incorrect unit: expected 0, got %d", len(unit2.Abilities))
	}

	if len(unit2.Keywords) != 0 {
		t.Errorf("service linked keyword to incorrect unit: expected 0, got %d", len(unit2.Keywords))
	}
}

func TestGetUnitByID_NotFound(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()
	randomID := uuid.New()

	_, err := GetUnitByID(s, ctx, randomID)
	if !errors.Is(err, appErr.ErrNotFound) {
		t.Fatalf("expected 'resource not found', got %v", err)
	}
}

func TestGetUnitByID_EmptyAssociations(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	unit, err := GetUnitByID(s, ctx, unitID)
	if err != nil {
		t.Fatalf("failed to get unit, %v", err)
	}

	if unit.Abilities == nil {
		t.Errorf("expected abilities slice to be initalized")
	}

	if len(unit.Abilities) != 0 {
		t.Errorf("expected 0 abilities, got %d", len(unit.Abilities))
	}

	if unit.Weapons == nil {
		t.Errorf("expected weapons slice to be intialized")
	}

	if len(unit.Weapons) != 0 {
		t.Errorf("expected 0 weapons, got %d", len(unit.Weapons))
	}

	if unit.Keywords == nil {
		t.Errorf("expected keywords to be initalized")
	}

	if len(unit.Keywords) != 0 {
		t.Errorf("expected 0 keywords, got %d", len(unit.Keywords))
	}
}
