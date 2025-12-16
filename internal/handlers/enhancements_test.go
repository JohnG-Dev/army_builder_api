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

func TestGetEnahancements_ReturnsEnhancements(t *testing.T) {
	s := setupTestDB(t)

	gameID := createTestGame(t, s)
	factionID := createTestFaction(t, s, gameID)
	createTestEnhancement(t, s, factionID)

	handler := &EnhancementsHandlers{S: s}

	req := httptest.NewRequest(http.MethodGet, "/enhancements/", nil)
	w := httptest.NewRecorder()

	handler.GetEnhancements(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test Enhancement") {
		t.Errorf("expected body to contain 'Test Enhancement', got %s", bodyStr)
	}
}

func TestGetEnhancements_FilterByFactionID(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	gameID := createTestGame(t, s)
	factionID1 := createTestFactionWithName(t, s, gameID, "StormCast")
	factionID2 := createTestFactionWithName(t, s, gameID, "Skaven")

	handler := &EnhancementsHandlers{S: s}

	_, err := s.DB.CreateEnhancement(ctx, database.CreateEnhancementParams{
		FactionID:       factionID1,
		Name:            "Test StormCast Enhancement",
		EnhancementType: "artefact",
		Description:     "Test Description",
		Points:          10,
		IsUnique:        true,
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create stormcast enhancement, %v", err)
	}

	_, err = s.DB.CreateEnhancement(ctx, database.CreateEnhancementParams{
		FactionID:       factionID2,
		Name:            "Test Skaven Enhancement",
		EnhancementType: "artefact",
		Description:     "Test Description",
		Points:          10,
		IsUnique:        true,
		Version:         "1.0",
		Source:          "Test Source",
	})
	if err != nil {
		t.Fatalf("failed to create skaven enhancement, %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/enhancements?faction_id="+factionID1.String(), nil)
	w := httptest.NewRecorder()

	handler.GetEnhancements(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Test StormCast Enhancement") {
		t.Errorf("expected body to contain 'Test StormCast Enhancement', got %s", bodyStr)
	}

	if strings.Contains(bodyStr, "Test Skaven Enhancement") {
		t.Errorf("expected body to NOT contain 'Test Skaven Enhancement', got %s", bodyStr)
	}
}
