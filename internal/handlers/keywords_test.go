package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	defer res.Body.Close()

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
	defer res.Body.Close()

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
