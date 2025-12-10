package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestListFactions_ReturnsFaction(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)

	createTestFaction(t, s, gameID)

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions", nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Test Faction") {
		t.Errorf("expected 'Test Faction' in response, got %s", string(body))
	}
}

func TestListFactions_Empty(t *testing.T) {
	s := setupTestDB(t)

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions", nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "[]" {
		t.Errorf("expected empty array [], got %s", string(body))
	}
}

func TestListFactions_FilterByGameID(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	createTestFaction(t, s, gameID)

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions?game_id="+gameID.String(), nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Faction") {
		t.Errorf("expected 'Test Faction', got %s", bodyStr)
	}
}

func TestGetFactionByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	handler := &FactionsHandlers{S: s}

	randomID := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/factions/"+randomID, nil)
	req.SetPathValue("id", randomID)
	w := httptest.NewRecorder()

	handler.GetFactionByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", res.StatusCode)
	}
}

func TestGetFactionByID_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)

	factionID := createTestFaction(t, s, gameID)

	handler := &FactionsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/factions/"+factionID.String(), nil)
	req.SetPathValue("id", factionID.String())
	w := httptest.NewRecorder()

	handler.GetFactionByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "Test Faction") {
		t.Errorf("expected 'Test Faction' in response, got %s", bodyStr)
	}
}

func TestGetFactionsByName(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	createTestFaction(t, s, gameID)

	handler := &FactionsHandlers{S: s}
	req := httptest.NewRequest(http.MethodGet, "/factions?name=Test%20Faction", nil)
	w := httptest.NewRecorder()

	handler.GetFactions(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Faction") {
		t.Errorf("expected body to contain 'Test Faction' got %s", bodyStr)
	}
}
