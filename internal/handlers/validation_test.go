package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
)

func TestValidateArmy_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	armyReq := models.ArmyValidationRequest{
		GameID:      gameID,
		FactionID:   factionID,
		PointsLimit: 2000,
		Units: []models.ArmyUnit{
			{UnitID: unitID, Quantity: 2},
		},
	}

	jsonData, err := json.Marshal(armyReq)
	if err != nil {
		t.Fatalf("failed to marshal army req: %v", err)
	}

	bodyReader := bytes.NewBuffer(jsonData)
	handler := &ValidationHandlers{S: s}

	req := httptest.NewRequest(http.MethodPost, "/validation/", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ValidateArmy(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	var resp models.ValidationResponse

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if !resp.IsValid {
		t.Errorf("expected army to be valid, got errors: %v", resp.Errors)
	}

	expectedPoints := 200
	if resp.TotalPoints != expectedPoints {
		t.Errorf("expected %d points, got %d points", expectedPoints, resp.TotalPoints)
	}

	if len(resp.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(resp.Errors), resp.Errors)
	}
}
