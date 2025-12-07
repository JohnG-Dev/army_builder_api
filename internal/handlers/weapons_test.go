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

func TestGetWeapons_ReturnsWeapons(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	game, err := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Test Game",
		Edition: "Test Edition",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	faction, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:  game.ID,
		Name:    "Test Faction",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create faction: %v", err)
	}

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       faction.ID,
		Name:            "Test unit",
		Move:            6,
		Health:          4,
		Save:            "4+",
		Ward:            "6+",
		Control:         1,
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     5,
		MaxUnitSize:     10,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        false,
		Version:         "1.0",
		Source:          "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create unit: %v", err)
	}

	_, err = s.DB.CreateWeapon(ctx, database.CreateWeaponParams{
		UnitID:  unit.ID,
		Name:    "Test AoS Weapon",
		Range:   "12\"",
		Attacks: "3",
		ToHit:   "3+",
		ToWound: "3+",
		Rend:    "-1",
		Damage:  "2",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create weapon: %v", err)
	}

	handler := &WeaponsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/weapons", nil)

	w := httptest.NewRecorder()

	handler.GetWeapons(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS Weapon") {
		t.Errorf("expected body to contain 'Test AoS Weapon', got %s", bodyStr)
	}
}

func TestGetWeapons_EmptyDB(t *testing.T) {
	s := setupTestDB(t)

	handler := &WeaponsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/weapons", nil)

	w := httptest.NewRecorder()

	handler.GetWeapons(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "[]") {
		t.Errorf("expected body to contain '[]', got %s", bodyStr)
	}

}

func TestGetWeapons_FilterByUnitID(t *testing.T) {
	s := setupTestDB(t)

	ctx := context.Background()

	game, err := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Test Game",
		Edition: "Test Edition",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	faction, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:  game.ID,
		Name:    "Test Faction",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create faction: %v", err)
	}

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       faction.ID,
		Name:            "Test Unit",
		Description:     "A test unit",
		Move:            6,
		Health:          3,
		Save:            "4+",
		Ward:            "6+",
		Control:         1,
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     5,
		MaxUnitSize:     10,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        false,
		Version:         "1.0",
		Source:          "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create unit: %v", err)
	}

	_, err = s.DB.CreateWeapon(ctx, database.CreateWeaponParams{
		UnitID:  unit.ID,
		Name:    "Test AoS Weapon",
		Range:   "12\"",
		Attacks: "3",
		ToHit:   "3+",
		ToWound: "3+",
		Rend:    "-1",
		Damage:  "2",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create weapon: %v", err)
	}

	handler := &WeaponsHandlers{S: s}
	req := httptest.NewRequest(http.MethodGet, "/weapons?unit_id="+unit.ID.String(), nil)

	w := httptest.NewRecorder()

	handler.GetWeapons(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS Weapon") {
		t.Errorf("expected body to contain 'Test AoS Weapon', got %s", bodyStr)
	}
}

func TestGetWeaponByID_Success(t *testing.T) {
	s := setupTestDB(t)

	ctx := context.Background()

	game, err := s.DB.CreateGame(ctx, database.CreateGameParams{
		Name:    "Test Game",
		Edition: "Test Edition",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	faction, err := s.DB.CreateFaction(ctx, database.CreateFactionParams{
		GameID:  game.ID,
		Name:    "Test Faction",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create faction: %v", err)
	}

	unit, err := s.DB.CreateUnit(ctx, database.CreateUnitParams{
		FactionID:       faction.ID,
		Name:            "Test Unit",
		Move:            6,
		Health:          3,
		Save:            "3+",
		Ward:            "6+",
		Control:         1,
		Points:          100,
		SummonCost:      "",
		Banishment:      "",
		MinUnitSize:     3,
		MaxUnitSize:     6,
		MatchedPlay:     true,
		IsManifestation: false,
		IsUnique:        false,
		Version:         "1.0",
		Source:          "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create unit: %v", err)
	}

	_, err = s.DB.CreateWeapon(ctx, database.CreateWeaponParams{
		UnitID:  unit.ID,
		Name:    "Test AoS Weapon",
		Range:   "10\"",
		Attacks: "5",
		ToHit:   "3+",
		ToWound: "3+",
		Rend:    "-1",
		Damage:  "2",
		Version: "1.0",
		Source:  "Test Source",
	})

	if err != nil {
		t.Fatalf("failed to create weapon: %v", err)
	}

	handler := &WeaponsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/weapons/"+unit.ID.String(), nil)
	req.SetPathValue("id", unit.ID.String())

	w := httptest.NewRecorder()

	handler.GetWeaponByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS Weapon") {
		t.Errorf("expected body to contain 'Test AoS Weapon', got %s", bodyStr)
	}
}

func TestGetWeaponByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	handler := &WeaponsHandlers{S: s}

	nonExistentID := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/weapons/"+nonExistentID, nil)
	req.SetPathValue("id", nonExistentID)
	w := httptest.NewRecorder()

	handler.GetWeaponByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}

}

func TestGetWeaponsByID_InvalidUUID(t *testing.T) {
	s := setupTestDB(t)

	handler := &WeaponsHandlers{S: s}

	invalidID := "not-a-valid-uuid"

	req := httptest.NewRequest(http.MethodGet, "/weapons/"+invalidID, nil)
	req.SetPathValue("id", invalidID)
	w := httptest.NewRecorder()

	handler.GetWeaponByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}

}
