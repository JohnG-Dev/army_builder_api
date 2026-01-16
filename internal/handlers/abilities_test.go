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
	defer func() { _ = res.Body.Close() }() 

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
		GameID:    uuid.NullUUID{},
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
		GameID:    uuid.NullUUID{},
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
	defer func() { _ = res.Body.Close() }() 

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

func TestGetAbilities_FilterByFactionID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID1 := createTestFactionWithName(t, s, gameID, "Stormcast")
	factionID2 := createTestFactionWithName(t, s, gameID, "Skaven")

	_, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    uuid.NullUUID{},
		FactionID: database.UUIDToNullUUID(factionID1),
		GameID:    uuid.NullUUID{},
		Name:      "Test Stormcast Ability",
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test stormcast ability, %v", err)
	}

	_, err = s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    uuid.NullUUID{},
		FactionID: database.UUIDToNullUUID(factionID2),
		GameID:    uuid.NullUUID{},
		Name:      "Test Skaven Ability",
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to creste test skaven ability, %v", err)
	}

	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities?faction_id="+factionID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetAbilities(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Stormcast Ability") {
		t.Errorf("expected body to contain 'test stormcast ability', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test Skaven Ability") {
		t.Errorf("expected body to NOT contain 'test skaven ability', got %s", bodyStr)
	}
}

func TestGetAbilites_FilterByType(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	_, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID),
		FactionID: uuid.NullUUID{},
		GameID:    uuid.NullUUID{},
		Name:      "Test Stormcast Prayer",
		Type:      "Prayer",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test stormcast prayer, %v", err)
	}

	_, err = s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID),
		FactionID: uuid.NullUUID{},
		GameID:    uuid.NullUUID{},
		Name:      "Test Stormcast Passive",
		Type:      "Passive",
		Phase:     "End Of Turn",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test stormcast passive, %v", err)
	}

	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities?type=Prayer", nil)
	w := httptest.NewRecorder()

	handler.GetAbilities(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Prayer") {
		t.Errorf("expected body to contain 'spell' got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Passive") {
		t.Errorf("expected body to NOT contain 'passive' got %s", bodyStr)
	}
}

func TestGetAbilities_FilterByPhase(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	_, err := s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID),
		FactionID: uuid.NullUUID{},
		GameID:    uuid.NullUUID{},
		Name:      "Test Stormcast ability",
		Type:      "Spell",
		Phase:     "Hero",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test stormcast abilitiy hero, %v", err)
	}

	_, err = s.DB.CreateAbility(ctx, database.CreateAbilityParams{
		UnitID:    database.UUIDToNullUUID(unitID),
		FactionID: uuid.NullUUID{},
		GameID:    uuid.NullUUID{},
		Name:      "Test StormCast Ability",
		Type:      "Passive",
		Phase:     "Shooting",
		Version:   "1.0",
		Source:    "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create test stormcast ability passive, %v", err)
	}

	handler := &AbilitiesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/abilities?phase=Hero", nil)
	w := httptest.NewRecorder()
	handler.GetAbilities(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Hero") {
		t.Errorf("exected body to contain 'hero', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Shooting") {
		t.Errorf("expected body to NOT contain 'shooting', got %s", bodyStr)
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
	defer func() { _ = res.Body.Close() }() 

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
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}
}
