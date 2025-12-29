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
		t.Fatalf("failed to get abilities 2: %v", err)
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

func TestGetAbilityByID_WithEffects(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)
	abilityID := createTestAbilityUnit(t, s, unitID)

	_, err := s.DB.CreateAbilityEffect(ctx, database.CreateAbilityEffectParams{
		AbilityID:   abilityID,
		Stat:        "Save",
		Modifier:    1,
		Condition:   "7+ prayer points",
		Description: "Bless Chosen Unit to Increase Save Score",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create ability effect save: %v", err)
	}

	_, err = s.DB.CreateAbilityEffect(ctx, database.CreateAbilityEffectParams{
		AbilityID:   abilityID,
		Stat:        "Rend",
		Modifier:    1,
		Condition:   "7+ prayer points",
		Description: "Bless Chosen Unit to Increase Rend",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create ability effect rend: %v", err)
	}

	ability, err := GetAbilityByID(s, ctx, abilityID)
	if err != nil {
		t.Fatalf("failted to get ability: %v", err)
	}

	if ability.Effects == nil {
		t.Fatalf("failed to initalize ability effects")
	}

	if ability.UnitID == nil || *ability.UnitID != unitID {
		t.Errorf("ability lost its unit id pointer during the fetch")
	}

	if len(ability.Effects) != 2 {
		t.Errorf("expected 2 ability effects, got %d", len(ability.Effects))
	}

	foundSave := false
	foundRend := false

	for _, eff := range ability.Effects {
		if eff.AbilityID != abilityID {
			t.Errorf("effect %v is linked to wrong ability: expected %v, got %v", eff.ID, abilityID, eff.AbilityID)
		}

		switch eff.Stat {
		case "Save":
			foundSave = true
			if eff.Modifier != 1 {
				t.Errorf("Save modifier: expected 1, got %d", eff.Modifier)
			}
		case "Rend":
			foundRend = true
			if eff.Modifier != 1 {
				t.Errorf("Rend modifier: expected 1, got %d", eff.Modifier)
			}
		default:
			t.Errorf("found unexpected stat in effects: %s", eff.Stat)
		}
	}

	if !foundSave {
		t.Errorf("did not find 'Save' effect in result")
	}

	if !foundRend {
		t.Errorf("did not find 'Rend' effect in result")
	}
}
