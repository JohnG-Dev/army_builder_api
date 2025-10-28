package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type AbilitiesHandlers struct {
	S *state.State
}

func (h *AbilitiesHandlers) GetAbilities(w http.ResponseWriter, r *http.Request) {
	abilities, err := services.GetAllAbilities(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "abilities not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch abilities", err)
		}
		logRequestError(h.S, r, "failed to fetch abiilities", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched abilities", zap.Int("count", len(abilities)))
	respondWithJSON(w, http.StatusOK, abilities)
}

func (h *AbilitiesHandlers) GetAbilitiesForUnit(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	factionID, err := uuid.Parse(factionIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid faction id", err)
		return
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
		return
	}

	logRequestInfo(h.S, r, "Successfuly fetched abilities", zap.Int("count", len(abilities)))
	respondWithJSON(w, http.StatusOK, abilities)
}

func (h *AbilitiesHandlers) GetAbilityByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing ability id", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ability id", err)
		return
	}

	ability, err := services.GetAbilityByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "ability not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch ability", err)
		}

		logRequestError(h.S, r, "failed to fetch ability", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched ability")
	respondWithJSON(w, http.StatusOK, ability)
}

func (h *AbilitiesHandlers) GetAbilityByType(w http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")

	if typeStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing ability type", nil)
		return
	}

	ability, err := services.GetAbilityByType(h.S, r.Context(), typeStr)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing ability type", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "ability not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch ability", err)
		}

		logRequestError(h.S, r, "failed to fetch ability", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched ability")
	respondWithJSON(w, http.StatusOK, ability)
}

func (h *AbilitiesHandlers) GetAbilityByPhase(w http.ResponseWriter, r *http.Request) {
	phaseStr := r.URL.Query().Get("phase")

	if phaseStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing phase type", nil)
		return
	}

	ability, err := services.GetAbilityByPhase(h.S, r.Context(), phaseStr)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing phase", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "ability not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch ability", err)
		}

		logRequestError(h.S, r, "failed to fetch ability", err)
	}

	logRequestInfo(h.S, r, "Successfully fetched ability")
	respondWithJSON(w, http.StatusOK, ability)
}
