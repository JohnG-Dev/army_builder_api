package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type AbilitiesHandlers struct {
	S *state.State
}

func (h *AbilitiesHandlers) GetAbilitiesForUnit(w http.ResponseWriter, r *http.Request) {
	unitIDStr := r.URL.Query().Get("unit_id")

	if unitIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing unit id", nil)
	}

	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
	}

	abilities, err := services.GetAbilitiesForUnit(h.S, r.Context(), unitID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingUnitID):
			respondWithError(w, http.StatusBadRequest, "unit id required", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "unit not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch abilities", err)
		}

		logRequestError(h.S, r, "failed to fetch abilities", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched abilities", zap.Int("count", len(abilities)))
	respondWithJSON(w, http.StatusOK, abilities)
}

func (h *AbilitiesHandlers) GetAbilitiesForFaction(w http.ResponseWriter, r *http.Request) {
	factionIDStr := r.URL.Query().Get("faction_id")

	if factionIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing faction id", nil)
	}

	factionID, err := uuid.Parse(factionIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid faction id", err)
	}

	abilities, err := services.GetAbilitiesForFaction(h.S, r.Context(), factionID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingFactionID):
			respondWithError(w, http.StatusBadRequest, "faction id required", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "abilities not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch abilities", err)
		}

		logRequestError(h.S, r, "failed to fetch abilities", err)
	}

	logRequestInfo(h.S, r, "Successfuly fetched abilities", zap.Int("count", len(abilities)))
	respondWithJSON(w, http.StatusOK, abilities)
}

func (h *AbilitiesHandlers) GetAbilityByID(w http.ResponseWriter, r *http.Request)
