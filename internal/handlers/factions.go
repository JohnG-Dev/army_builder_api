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

type FactionsHandlers struct {
	S *state.State
}

func (h *FactionsHandlers) GetFactions(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name != "" {
		h.getFactionsByName(w, r)
		return
	}

	h.getAllFactions(w, r)
}

func (h *FactionsHandlers) getAllFactions(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.URL.Query().Get("game_id")
	isAoRStr := r.URL.Query().Get("is_army_of_renown")
	isRoRStr := r.URL.Query().Get("is_regiment_of_renown")
	parentIDStr := r.URL.Query().Get("parent_id")

	filter := services.FactionFilter{}

	if gameIDStr != "" {
		id, err := uuid.Parse(gameIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid game id", err)
			return
		}
		filter.GameID = &id
	}

	if isAoRStr != "" {
		val := isAoRStr == "true"
		filter.IsArmyOfRenown = &val
	}

	if isRoRStr != "" {
		val := isRoRStr == "true"
		filter.IsRegimentOfRenown = &val
	}

	if parentIDStr != "" {
		id, err := uuid.Parse(parentIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid parent id", err)
			return
		}
		filter.ParentFactionID = &id
	}

	factions, err := services.GetFactions(h.S, r.Context(), filter)
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

func (h *FactionsHandlers) getFactionsByName(w http.ResponseWriter, r *http.Request) {
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
