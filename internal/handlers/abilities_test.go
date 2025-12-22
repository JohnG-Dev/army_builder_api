package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

func TestGetAbilities_ReturnsAbilities(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFactionWithName(t, s, gameID, "Test Stormcast")
	unitID := createTestUnitWithName(t, s, factionID, "Test Stormcast Priest")

	createTestAbilityFaction(t, s, factionID)
	createTestAbilityUnit(t, s, unitID)

	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities/", nil)
	w := httptest.NewRecorder()

	handler.GetAbilities(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Ability") {
		t.Errorf("expected body to contain 'test ability', got %s", bodyStr)
	}
}

func TestGetAbilities_FilterByUnitID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID1 := createTestUnit(t, s, factionID)
	unitID2 := createTestUnit(t, s, factionID)
	_, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID1),
		FactionID: uuid.NullUUID{},
		Name:      "Test Wizard Ability",
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test wizard ability, %v", err)
	}
	_, err = s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID2),
		FactionID: uuid.NullUUID{},
		Name:      "Test Melee Ability",
		Type:      "Combat",
		Phase:     "Charge",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test melee ability, %v", err)
	}

	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities?unit_id="+unitID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetAbilities(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Wizard Ability") {
		t.Errorf("expected body to contain 'test wizard ability', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test Melee Ability") {
		t.Errorf("expected body to NOT contain 'test melee ability' got %s", bodyStr)
	}
}

func TestGetAbilityByID_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)
	abilityID := createTestAbilityUnit(t, s, unitID)

	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities/"+abilityID.String(), nil)
	req.SetPathValue("id", abilityID.String())
	w := httptest.NewRecorder()

	handler.GetAbilityByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status cde 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Ability") {
		t.Errorf("expected body to contain 'test ability', got %s", bodyStr)
	}
}

func TestGetAbilityByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	randomID := uuid.New()
	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities/"+randomID.String(), nil)
	req.SetPathValue("id", randomID.String())
	w := httptest.NewRecorder()

	handler.GetAbilityByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}
}
