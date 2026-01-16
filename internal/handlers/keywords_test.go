package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
)

func TestGetKeywords_ReturnsKeywords(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	createTestKeyword(t, s, gameID)

	handler := &KeywordsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/keywords/", nil)
	w := httptest.NewRecorder()

	handler.GetKeywords(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Keyword") {
		t.Errorf("expected body to contain 'Test Keyword', got %s", bodyStr)
	}
}

func TestGetKeywords_FilterByGameID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID1 := createTestGameWithName(t, s, "Age of Sigmar")
	gameID2 := createTestGameWithName(t, s, "Warhammer 40k")

	_, err := s.DB.CreateKeyword(ctx, database.CreateKeywordParams{
		GameID:      gameID1,
		Name:        "Test AoS Keyword",
		Description: "Test Keyword Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create AoS keyword, %v", err)
	}

	_, err = s.DB.CreateKeyword(ctx, database.CreateKeywordParams{
		GameID:      gameID2,
		Name:        "Test Warhammer 40k Keyword",
		Description: "Test Keyword Description",
		Version:     "1.0",
		Source:      "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create 40k keyword, %v", err)
	}

	handler := &KeywordsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/keywords?game_id="+gameID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetKeywords(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test AoS Keyword") {
		t.Errorf("expected body to contain 'Test AoS Keyword', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test Warhammer 40k Keyword") {
		t.Errorf("expected body to NOT contain 'Test Warhammer 40k Keyword', got %s", bodyStr)
	}
}

func TestGetKeywords_FilterByUnitID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID1 := createTestUnitWithName(t, s, factionID, "Judicators")
	unitID2 := createTestUnitWithName(t, s, factionID, "Liberators")
	keywordID1 := createTestKeywordWithName(t, s, gameID, "Judicators Keyword")
	keywordID2 := createTestKeywordWithName(t, s, gameID, "Liberators Keyword")

	_ = s.DB.AddKeywordToUnit(ctx, database.AddKeywordToUnitParams{
		UnitID:    unitID1,
		KeywordID: keywordID1,
		Value:     "Test Judicator Value",
	})

	_ = s.DB.AddKeywordToUnit(ctx, database.AddKeywordToUnitParams{
		UnitID:    unitID2,
		KeywordID: keywordID2,
		Value:     "Test Liberators Value",
	})

	handler := &KeywordsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/keywords?unit_id="+unitID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetKeywords(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Judicator Value") {
		t.Errorf("expected body to contain 'test judicators value', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test Liberators Value") {
		t.Errorf("expected body to NOT contain 'test liberators value', got %s", bodyStr)
	}
}

func TestGetUnitsWithKeyword_Success(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID := createTestUnitWithName(t, s, factionID, "Test Judicator")
	keywordID := createTestKeyword(t, s, gameID)
	_ = s.DB.AddKeywordToUnit(ctx, database.AddKeywordToUnitParams{
		UnitID:    unitID,
		KeywordID: keywordID,
		Value:     "Judicator",
	})

	handler := &KeywordsHandlers{S: s}

	keywordName := "Test Keyword"
	escapedName := url.PathEscape(keywordName)
	req := httptest.NewRequest(http.MethodGet, "/keywords/"+escapedName+"/units", nil)
	req.SetPathValue("name", keywordName)

	w := httptest.NewRecorder()
	handler.GetUnitsWithKeyword(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Judicator") {
		t.Errorf("expected body to contain 'test judicator', got %s", bodyStr)
	}
}

func TestGetUnitsWithKeywordAndValue(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	unitID1 := createTestUnitWithName(t, s, factionID, "Test Judicator")
	unitID2 := createTestUnitWithName(t, s, factionID, "Vanguard Raptors")
	keywordID := createTestKeyword(t, s, gameID)
	_ = s.DB.AddKeywordToUnit(ctx, database.AddKeywordToUnitParams{
		UnitID:    unitID1,
		KeywordID: keywordID,
		Value:     "Melee",
	})

	_ = s.DB.AddKeywordToUnit(ctx, database.AddKeywordToUnitParams{
		UnitID:    unitID2,
		KeywordID: keywordID,
		Value:     "Ranged",
	})

	handler := &KeywordsHandlers{S: s}

	keywordName := "Test Keyword"
	searchValue := "Melee"
	escapedName := url.PathEscape(keywordName)
	escapedValue := url.PathEscape(searchValue)

	req := httptest.NewRequest(http.MethodGet, "/keywords/"+escapedName+"/units/value/"+escapedValue, nil)
	req.SetPathValue("name", keywordName)
	req.SetPathValue("value", searchValue)

	w := httptest.NewRecorder()

	handler.GetUnitsWithKeywordAndValue(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Judicator") {
		t.Errorf("expected body to contain 'test judicator', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Vanguard Raptors") {
		t.Errorf("expected body to NOT contain 'vanguard raptors', got %s", bodyStr)
	}
}

func TestGetKeywordByID_Success(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	keywordID := createTestKeyword(t, s, gameID)

	handler := &KeywordsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/keywords/"+keywordID.String(), nil)
	req.SetPathValue("id", keywordID.String())

	w := httptest.NewRecorder()

	handler.GetKeywordByID(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Keyword") {
		t.Errorf("expected body to contain 'test keyword', got %s", bodyStr)
	}
}

func TestGetKeywordByID_NotFound(t *testing.T) {
	s := setupTestDB(t)

	randomID := uuid.New()

	handler := &KeywordsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/keywords/"+randomID.String(), nil)
	req.SetPathValue("id", randomID.String())

	w := httptest.NewRecorder()

	handler.GetKeywordByID(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }() 

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	}
}
