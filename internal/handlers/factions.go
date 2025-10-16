package handlers

import (
	"database/sql"
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
	gameIDStr := r.URL.Query().Get("game_id")

	var factions []database.Faction
	var err error

	if gameIDStr == "" {
		factions, err = services.GetFactions(h.S, r.Context(), nil)
		if err != nil {

			logRequestError(h.S, r, "failed to fetch factions", err)
			respondWithError(w, http.StatusInternalServerError, "failed to fetch factions", err)
			return
		}

		respondWithJSON(w, http.StatusOK, factions)
		return
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse uuid", err)
		return
	}

	factions, err := services.GetFactions(h.S, r.Context(), &gameID)
	if err != nil {
		logRequestError(h.S, r, "failed to fetch factions", err)

		respondWithError(w, http.StatusInternalServerError, "failed to find factions", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched Factions",
		zap.Int("count", len(factions)),
	)

	respondWithJSON(w, http.StatusOK, factions)
}

func (h *FactionsHandlers) GetFactionByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	factionID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse uuid", err)
		return
	}

	faction, err := services.GetFactionByID(h.S, r.Context(), factionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestError(h.S, r, "faction not found", err)

			respondWithError(w, http.StatusNotFound, "faction not found", nil)
		} else {
			logRequestError(h.S, r, "failed to fetch faction", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch faction", err)
		}

		return
	}

	logRequestInfo(h.S, r, "Successfully fetched faction")

	respondWithJSON(w, http.StatusOK, faction)
}
