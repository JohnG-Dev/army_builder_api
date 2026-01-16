package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGetWeapons_ReturnsWeapons(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	createTestWeapon(t, s, unitID)

	handler := &WeaponsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/weapons", nil)

	w := httptest.NewRecorder()

	handler.GetWeapons(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

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
	defer func() { _ = res.Body.Close() }() 

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

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)
	createTestWeapon(t, s, unitID)

	handler := &WeaponsHandlers{S: s}
	req := httptest.NewRequest(http.MethodGet, "/weapons?unit_id="+unitID.String(), nil)

	w := httptest.NewRecorder()

	handler.GetWeapons(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

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

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)
	weaponID := createTestWeapon(t, s, unitID)

	handler := &WeaponsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/weapons/"+weaponID.String(), nil)
	req.SetPathValue("id", weaponID.String())

	w := httptest.NewRecorder()

	handler.GetWeaponByID(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

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

	nonExistentID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/weapons/"+nonExistentID.String(), nil)
	req.SetPathValue("id", nonExistentID.String())
	w := httptest.NewRecorder()

	handler.GetWeaponByID(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

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
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", res.StatusCode)
	}
}
