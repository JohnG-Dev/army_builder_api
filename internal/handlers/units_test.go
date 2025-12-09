package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGetUnits_ReturnsUnits(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	createTestUnit(t, s, factionID)

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()

	handler.GetNonManifestationUnits(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Test Unit") {
		t.Errorf("expected 'Test Unit' in response, got %s", string(body))
	}
}

func TestGetUnits_EmptyDB(t *testing.T) {
	s := setupTestDB(t)

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units", nil)
	w := httptest.NewRecorder()

	handler.GetNonManifestationUnits(w, req)

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

func TestGetUnits_FilterByFactionID(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	createTestUnit(t, s, factionID)

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units?faction_id="+factionID.String(), nil)
	w := httptest.NewRecorder()

	handler.GetUnits(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Unit") {
		t.Errorf("expected 'Test Unit', got %s", bodyStr)
	}
}

func TestGetManifestations_ReturnsManifestations(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	createTestManifestation(t, s, factionID)

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/manifestations", nil)
	w := httptest.NewRecorder()

	handler.GetManifestations(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Test Manifestation") {
		t.Errorf("expected 'Test Manifestation' in response, got %s", string(body))
	}
}

func TestGetUnitByID_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnit(t, s, factionID)

	handler := &UnitsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/units/"+unitID.String(), nil)
	req.SetPathValue("id", unitID.String())
	w := httptest.NewRecorder()

	handler.GetUnitByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Test Unit") {
		t.Errorf("expected 'Test Unit' in response, got %s", string(body))
	}
}

func TestGetUnitByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	handler := &UnitsHandlers{S: s}

	randomID := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/units/"+randomID, nil)
	req.SetPathValue("id", randomID)
	w := httptest.NewRecorder()

	handler.GetUnitByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", res.StatusCode)
	}
}
