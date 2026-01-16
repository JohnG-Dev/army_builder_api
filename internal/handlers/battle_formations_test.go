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

func TestGetBattleFormations_ReturnsBattleFormations(t *testing.T) {
	s := setupTestDB(t)

	handler := &BattleFormationsHandlers{S: s}

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	createTestBattleFormation(t, s, gameID, factionID)

	req := httptest.NewRequest(http.MethodGet, "/battle_formations/", nil)
	w := httptest.NewRecorder()

	handler.GetBattleFormations(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test BattleFormation") {
		t.Errorf("expected body to contain 'Test BattleFormation', got %s", bodyStr)
	}
}

func TestGetBattleFormations_FilterByGameID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	handler := &BattleFormationsHandlers{S: s}

	gameID1 := createTestGameWithName(t, s, "Age of Sigmar")
	gameID2 := createTestGameWithName(t, s, "Warhammer 40k")
	factionID1 := createTestFaction(t, s, gameID1)
	factionID2 := createTestFaction(t, s, gameID2)

	_, err := s.DB.CreateBattleFormation(ctx, database.CreateBattleFormationParams{
		GameID:      gameID1,
		FactionID:   factionID1,
		Name:        "Test AoS BattleFormation",
		Description: "AoS Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create battle formation, %v", err)
	}

	_, err = s.DB.CreateBattleFormation(ctx, database.CreateBattleFormationParams{
		GameID:      gameID2,
		FactionID:   factionID2,
		Name:        "Test 40k BattleFormation",
		Description: "40k Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create battle formation, %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/battle_formations?game_id="+gameID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetBattleFormations(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS BattleFormation") {
		t.Errorf("expected body to contain 'Test AoS BattleFormation', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test 40k BattleFormation") {
		t.Errorf("expected body to NOT contain 'Test 40k BattleFormation', got %s", bodyStr)
	}
}

func TestBattleFormations_FilterByFactionID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	handler := &BattleFormationsHandlers{S: s}

	gameID := createTestGame(t, s)
	factionID1 := createTestFactionWithName(t, s, gameID, "StormCast")
	factionID2 := createTestFactionWithName(t, s, gameID, "Tzeench")

	_, err := s.DB.CreateBattleFormation(ctx, database.CreateBattleFormationParams{
		GameID:      gameID,
		FactionID:   factionID1,
		Name:        "StormCast BattleFormation",
		Description: "StormCast Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create stormcast battle formation, %v", err)
	}

	_, err = s.DB.CreateBattleFormation(ctx, database.CreateBattleFormationParams{
		GameID:      gameID,
		FactionID:   factionID2,
		Name:        "Tzeench BattleFormation",
		Description: "Tzeench Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create tzeench battle formation, %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/battle_formations?faction_id="+factionID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetBattleFormations(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "StormCast BattleFormation") {
		t.Errorf("exected body to contain 'StormCast BattleFormation', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Tzeench BattleFormation") {
		t.Errorf("exected body to NOT contain 'Tzeench BattleFormation', got %s", bodyStr)
	}
}

func TestGetBattleFormationByID_Success(t *testing.T) {
	s := setupTestDB(t)

	handler := &BattleFormationsHandlers{S: s}

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	battleFormationID := createTestBattleFormation(t, s, gameID, factionID)

	req := httptest.NewRequest(http.MethodGet, "/battle_formations/"+battleFormationID.String(), nil)
	req.SetPathValue("id", battleFormationID.String())
	w := httptest.NewRecorder()

	handler.GetBattleFormationByID(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test BattleFormation") {
		t.Errorf("expected body to contain 'Test BattleFormation', got %s", bodyStr)
	}
}

func TestGetBattleFormationByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	handler := &BattleFormationsHandlers{S: s}

	randomID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/battle_formations/"+randomID.String(), nil)
	req.SetPathValue("id", randomID.String())
	w := httptest.NewRecorder()

	handler.GetBattleFormationByID(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}
}
