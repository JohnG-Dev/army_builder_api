package services

import (
	"context"
	"testing"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
)

func TestValidateArmy_Logic(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	myFaction := createTestFaction(t, s, gameID)
	enemyFaction := createTestFactionWithName(t, s, gameID, "Other")

	unit := createTestUnit(t, s, myFaction)
	uniqueUnit := createTestUniqueUnit(t, s, myFaction)
	enemyUnit := createTestUnit(t, s, enemyFaction)

	tests := []struct {
		name          string
		req           models.ArmyValidationRequest
		expectedValid bool
		expectedCount int
	}{
		{
			name: "Valid Army",
			req: models.ArmyValidationRequest{
				FactionID:   myFaction,
				PointsLimit: 2000,
				Units:       []models.ArmyUnit{{UnitID: unit, Quantity: 5}},
			},
			expectedValid: true,
			expectedCount: 0,
		},
		{
			name: "Unique Violation",
			req: models.ArmyValidationRequest{
				FactionID:   myFaction,
				PointsLimit: 2000,
				Units:       []models.ArmyUnit{{UnitID: uniqueUnit, Quantity: 2}},
			},
			expectedValid: false,
			expectedCount: 2,
		},
		{
			name: "Points Violation",
			req: models.ArmyValidationRequest{
				FactionID:   myFaction,
				PointsLimit: 500,
				Units:       []models.ArmyUnit{{UnitID: unit, Quantity: 8}},
			},
			expectedValid: false,
			expectedCount: 1,
		},
		{
			name: "Invalid Unit",
			req: models.ArmyValidationRequest{
				FactionID:   myFaction,
				PointsLimit: 2000,
				Units:       []models.ArmyUnit{{UnitID: enemyUnit, Quantity: 4}},
			},
			expectedValid: false,
			expectedCount: 1,
		},
		{
			name: "Size Violation (Below Min Limit)",
			req: models.ArmyValidationRequest{
				FactionID:   myFaction,
				PointsLimit: 2000,
				Units:       []models.ArmyUnit{{UnitID: unit, Quantity: 1}},
			},
			expectedValid: false,
			expectedCount: 1,
		},
		{
			name: "Size Violation (Above max Limit)",
			req: models.ArmyValidationRequest{
				FactionID:   myFaction,
				PointsLimit: 2000,
				Units:       []models.ArmyUnit{{UnitID: unit, Quantity: 9}},
			},
			expectedValid: false,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ValidateArmy(s, ctx, tt.req)
			if err != nil {
				t.Fatalf("unexpected system error: %v", err)
			}

			if resp.IsValid != tt.expectedValid {
				t.Errorf("expected IsValid to be %v, got %v", tt.expectedValid, resp.IsValid)
			}

			if len(resp.Errors) != tt.expectedCount {
				t.Errorf("expected %d errors, got %d: %v", tt.expectedCount, len(resp.Errors), resp.Errors)
			}
		})
	}
}
