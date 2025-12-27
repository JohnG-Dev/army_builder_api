package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetFactionsByName(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID1 := createTestFactionWithName(t, s, gameID, "Stormcast Eternals")
	factionID2 := createTestFactionWithName(t, s, gameID, "Stormcast Vanguard")
	factionID3 := createTestFactionWithName(t, s, gameID, "Skaven")

	factions, err := GetFactionsByName(s, ctx, "Storm")
	if err != nil {
		t.Fatalf("failed to get factions by name: %v", err)
	}

	if len(factions) != 2 {
		t.Fatalf("expected 2 factions, got %d", len(factions))
	}

	foundIDs := make(map[uuid.UUID]bool)
	for _, f := range factions {
		foundIDs[f.ID] = true
	}

	if !foundIDs[factionID1] {
		t.Errorf("missing expected faction: 'Stormcast Eternals'")
	}

	if !foundIDs[factionID2] {
		t.Errorf("missing expected faction: 'Stormcast Vanguard'")
	}

	if foundIDs[factionID3] {
		t.Errorf("unwanted faction included in results: 'Skaven'")
	}
}
