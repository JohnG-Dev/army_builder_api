package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
			{UnitID: unitID, Quantity: 4},
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
	defer func() { _ = res.Body.Close() }() 

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

	expectedPoints := 400
	if resp.TotalPoints != expectedPoints {
		t.Errorf("expected %d points, got %d points", expectedPoints, resp.TotalPoints)
	}

	if len(resp.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(resp.Errors), resp.Errors)
	}
}

func TestValidateArmy_OverPoints(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	armyReq := models.ArmyValidationRequest{
		GameID:      gameID,
		FactionID:   factionID,
		PointsLimit: 2000,
		Units: []models.ArmyUnit{
			{UnitID: unitID, Quantity: 21},
		},
	}

	jsonData, err := json.Marshal(armyReq)
	if err != nil {
		t.Fatalf("falied to marshal army req: %v", err)
	}

	bodyReader := bytes.NewBuffer(jsonData)
	handler := &ValidationHandlers{S: s}

	req := httptest.NewRequest(http.MethodPost, "/validation/", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ValidateArmy(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	var resp models.ValidationResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if resp.TotalPoints != 2100 {
		t.Errorf("expected 2100 points, got %d", resp.TotalPoints)
	}

	if resp.IsValid {
		t.Errorf("expected valid to be is false because points should be over the limit")
	}

	if len(resp.Errors) == 0 {
		t.Errorf("expected error messagesbut got none")
	}
}

func TestValidateArmy_WrongFaction(t *testing.T) {
	s := setupTestDB(t)
	gameID := createTestGame(t, s)
	factionID1 := createTestFactionWithName(t, s, gameID, "Stormcast")
	factionID2 := createTestFactionWithName(t, s, gameID, "Tzeench")
	unitID := createTestUnit(t, s, factionID2)

	armyReq := models.ArmyValidationRequest{
		GameID:      gameID,
		FactionID:   factionID1,
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
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	var resp models.ValidationResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if resp.IsValid {
		t.Errorf("expected is valid to be false because army faction is invalid for unit")
	}

	if len(resp.Errors) == 0 {
		t.Errorf("expected error messages but got none")
	}

	allErrors := strings.Join(resp.Errors, " ")

	if !strings.Contains(allErrors, "faction") {
		t.Errorf("expected error to contain a message about faction mismatch, got %v", resp.Errors)
	}
}
