package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRules_ReturnsRules(t *testing.T) {
	s := setupTestDB(t)

	handler := &RulesHandlers{S: s}

	gameID := createTestGame(t, s)
	createTestRule(t, s, gameID)

	req := httptest.NewRequest(http.MethodGet, "/rules/", nil)
	w := httptest.NewRecorder()

	handler.GetRules(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Rule") {
		t.Errorf("expected body to contain 'Test Rule', got %s", bodyStr)
	}
}

func TestGetRules_EmptyDB(t *testing.T) {
	s := setupTestDB(t)

	handler := &RulesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/rules/", nil)
	w := httptest.NewRecorder()

	handler.GetRules(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "[]") {
		t.Errorf("expected body to contain '[]', got %s", bodyStr)
	}
}

func TestGetRules_FilterByGameID(t *testing.T) {
	s := setupTestDB(t)

	handler := &RulesHandlers{S: s}
	gameID := createTestGame(t, s)
	createTestRule(t, s, gameID)

	req := httptest.NewRequest(http.MethodGet, "/rules?game_id="+gameID.String(), nil)
	req.SetPathValue("id", gameID.String())
	w := httptest.NewRecorder()

	handler.GetRuleByID(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Rule") {
		t.Errorf("expected body to contain 'Test Rule', got %s", bodyStr)
	}
}
