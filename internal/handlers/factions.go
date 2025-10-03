package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FactionsHandlers struct {
	S *state.State
}

func (h *FactionsHandlers) ListFactions(w http.ResponseWriter, r *http.Request) {

	gameIDString := r.URL.Query().Get("game_id")

	if gameIDString == "" {
		dbFactionList, err := services.ListFactions(h.S, r.Context(), nil)
		if err != nil {

			logRequestError(h.S, r, "failed to fetch factions", err)
			respondWithError(w, http.StatusInternalServerError, "failed to fetch factions", err)
			return
		}

		respondWithJSON(w, http.StatusOK, dbFactionList)
		return
	}

	gameID, err := uuid.Parse(gameIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse uuid", err)
		return
	}

	dbFactionList, err := services.ListFactions(h.S, r.Context(), &gameID)
	if err != nil {
		logRequestError(h.S, r, "failed to fetch factions", err)

		respondWithError(w, http.StatusInternalServerError, "failed to find factions", err)
		return
	}

	logRequestInfo(h.S, r, "Fetched Factions successfully",
		zap.Int("count", len(dbFactionList)),
	)

	respondWithJSON(w, http.StatusOK, dbFactionList)
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
			logRequestInfo(h.S, r, "Faction not found")

			respondWithError(w, http.StatusNotFound, "Faction not found", nil)
		} else {
			logRequestError(h.S, r, "failed to fetch faction", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch faction", err)
		}

		return
	}

	logRequestInfo(h.S, r, "Fetched faction successfully")

	respondWithJSON(w, http.StatusOK, faction)
}
