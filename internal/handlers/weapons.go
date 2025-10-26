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

type WeaponsHandlers struct {
	S *state.State
}

func (h *WeaponsHandlers) GetWeaponsForUnit(w http.ResponseWriter, r *http.Request) {
	unitIDStr := r.URL.Query().Get("unit_id")

	if unitIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing unit id", nil)
		return
	}

	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
		return
	}

	weapons, err := services.GetWeaponsForUnit(h.S, r.Context(), &unitID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingUnitID):
			respondWithError(w, http.StatusBadRequest, "unit id required", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "weapons not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch weapons", err)
		}

		logRequestError(h.S, r, "failed to fetch weapons", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched weapons", zap.Int("count", len(weapons)))
	respondWithJSON(w, http.StatusOK, weapons)
}

func (h *WeaponsHandlers) GetWeaponByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid weapon id", err)
		return
	}

	weapon, err := services.GetWeaponByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "weapon not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch weapon", err)
		}

		logRequestError(h.S, r, "failed to fetch weapon", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched weapon")
	respondWithJSON(w, http.StatusOK, weapon)

}
