package services

import (
	"context"
	"testing"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

func TestGetAbilitiesForUnit(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID1 := createTestUnitWithName(t, s, factionID, "Judicator")
	unitID2 := createTestUnitWithName(t, s, factionID, "Liberator")
	createTestAbilityUnitWithName(t, s, unitID1, "Bless")
	createTestAbilityUnitWithName(t, s, unitID1, "Extremis Chamber")

	abilities1, err := GetAbilitiesForUnit(s, ctx, unitID1)
	if err != nil {
		t.Fatalf("failed to get abilties 1: %v", err)
	}
	abilities2, err := GetAbilitiesForUnit(s, ctx, unitID2)
	if err != nil {
		t.Fatalf("failed to get abilities 2: %v, err")
	}

	if len(abilities1) != 2 {
		t.Errorf("expected 2 abilities, got %d", len(abilities1))
	}

	if abilities1[0].UnitID == nil || *abilities1[0].UnitID != unitID1 {
		t.Errorf("incorrect unit mapping")
	}

	if abilities1[0].FactionID != nil {
		t.Errorf("expected faction ID to be nil for abilties")
	}

	if len(abilities2) != 0 {
		t.Errorf("isolation failure: expected 0 abilities in abilities '2', got %d", len(abilities2))
	}
}
