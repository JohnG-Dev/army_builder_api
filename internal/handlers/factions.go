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

type FactionsHandlers struct {
	S *state.State
}

func (h *FactionsHandlers) GetFactions(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.URL.Query().Get("game_id")
	var factions []models.Faction
	var err error

	if gameIDStr == "" {
		factions, err = services.GetFactions(h.S, r.Context(), nil)
	} else {
		gameID, parseErr := uuid.Parse(gameIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid game id", parseErr)
			return
		}
		factions, err = services.GetFactions(h.S, r.Context(), &gameID)
	}

	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingFactionID), errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing or invalid id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "factions not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch factions", err)
		}
		logRequestError(h.S, r, "failed to fetch factions", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched factions", zap.Int("count", len(factions)))
	respondWithJSON(w, http.StatusOK, factions)
}

func (h *FactionsHandlers) GetFactionByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	factionID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid faction id", err)
		return
	}

	faction, err := services.GetFactionByID(h.S, r.Context(), factionID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "faction not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch faction", err)
		}

		logRequestError(h.S, r, "failed to fetch faction", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched faction")
	respondWithJSON(w, http.StatusOK, faction)
}

func (h *FactionsHandlers) GetFactionsByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		respondWithError(w, http.StatusBadRequest, "missing name", nil)
		return
	}

	factions, err := services.GetFactionsByName(h.S, r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "factions not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch factions", err)
		}

		logRequestError(h.S, r, "failed to fetch factions", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched factions", zap.Int("count", len(factions)))
	respondWithJSON(w, http.StatusOK, factions)
}
