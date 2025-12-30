package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type ValidationHandlers struct {
	S *state.State
}

type ArmyValidationRequest struct {
	GameID      uuid.UUID  `json:"game_id"`
	FactionID   uuid.UUID  `json:"faction_id"`
	PointsLimit int        `json:"points_limit"`
	Units       []ArmyUnit `json:"units"`
}

type ArmyUnit struct {
	UnitID   uuid.UUID `json:"unit_id"`
	Quantity int       `json:"quantity"`
}

func (h *ValidationHandlers) ValidateArmy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ArmyValidationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := services.ValidateArmy(h.S, r.Context(), req)
	if err != nil {
		logRequestError(h.S, r, "validation service failure", err)
		respondWithError(w, http.StatusInternalServerError, "failed to validate army", err)
	}

	logRequestInfo(h.S, r, "Army validation completed", zap.Bool("is_valid", resp.IsValid))
	respondWithJSON(w, http.StatusOK, resp)
}
