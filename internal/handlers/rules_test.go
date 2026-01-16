package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/google/uuid"
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
	defer func() { _ = res.Body.Close() }() 

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
	defer func() { _ = res.Body.Close() }() 

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
	ctx := context.Background()

	handler := &RulesHandlers{S: s}
	gameID1 := createTestGameWithName(t, s, "Age of Sigmar")
	gameID2 := createTestGameWithName(t, s, "Warhammer 40k")

	_, err := s.DB.CreateRule(ctx, database.CreateRuleParams{
		GameID:      gameID1,
		Name:        "Test AoS Rule",
		RuleType:    "core",
		Text:        "Rule text content",
		Description: "Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create AoS Rule: %v", err)
	}
	_, err = s.DB.CreateRule(ctx, database.CreateRuleParams{
		GameID:      gameID2,
		Name:        "Test Warhammer 40k Rule",
		RuleType:    "core",
		Text:        "Rule text content",
		Description: "Test Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create Test Warhammer 40k rule: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/rules?game_id="+gameID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetRules(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS Rule") {
		t.Errorf("expected body to contain 'Test AoS Rule', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test Warhammer 40k Rule") {
		t.Errorf("expected body not to contain 'Test Warhammer 40k Rule', got %s", bodyStr)
	}
}

func TestGetRules_FilterByType(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	handler := &RulesHandlers{S: s}

	gameID1 := createTestGameWithName(t, s, "Age of Sigmar")
	gameID2 := createTestGameWithName(t, s, "Warhammer 40k")

	_, err := s.DB.CreateRule(ctx, database.CreateRuleParams{
		GameID:      gameID1,
		Name:        "Test AoS Rule Core",
		RuleType:    "core",
		Text:        "Rule Text Content",
		Description: "Rule Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create rule, %v", err)
	}

	_, err = s.DB.CreateRule(ctx, database.CreateRuleParams{
		GameID:      gameID1,
		Name:        "Test AoS Rule Battle Tactic",
		RuleType:    "battle_tactic",
		Text:        "Rule Text Content",
		Description: "Rule Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create rule, %v", err)
	}

	_, err = s.DB.CreateRule(ctx, database.CreateRuleParams{
		GameID:      gameID2,
		Name:        "Test 40k Rule",
		RuleType:    "core",
		Text:        "Rule Text Content",
		Description: "Rule Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create rule, %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/rules?game_id="+gameID1.String()+"&type=core", nil)
	w := httptest.NewRecorder()

	handler.GetRules(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS Rule Core") {
		t.Errorf("expected body to contain 'Test AoS Rule Core', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test AoS Rule Battle Tactic") {
		t.Errorf("expected body to not contain 'Test AoS Rule Battle Tactic', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test 40k Rule") {
		t.Errorf("expected body to not contain 'Test 40k Rule', got %s", bodyStr)
	}
}

func TestGetRuleByID_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	ruleID := createTestRule(t, s, gameID)

	handler := &RulesHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/rules/"+ruleID.String(), nil)
	req.SetPathValue("id", ruleID.String())

	w := httptest.NewRecorder()

	handler.GetRuleByID(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Rule") {
		t.Errorf("expected body to contain 'Test Rule', got %s", bodyStr)
	}
}

func TestGetRuleByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	handler := &RulesHandlers{S: s}
	randomID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/rules/"+randomID.String(), nil)
	req.SetPathValue("id", randomID.String())

	w := httptest.NewRecorder()
	handler.GetRuleByID(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}
}
