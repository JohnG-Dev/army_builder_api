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

func ListFactions(s *state.State, w http.ResponseWriter, r *http.Request) {

	gameIDString := r.URL.Query().Get("game_id")

	if gameIDString == "" {
		dbFactionList, err := services.ListFactions(s, r.Context(), nil)
		if err != nil {

			logRequestError(s, r, "failed to fetch factions", err)
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

	dbFactionList, err := services.ListFactions(s, r.Context(), &gameID)
	if err != nil {
		logRequestError(s, r, "failed to fetch factions", err)

		respondWithError(w, http.StatusInternalServerError, "failed to find factions", err)
		return
	}

	logRequestInfo(s, r, "Fetched Factions successfully",
		zap.Int("count", len(dbFactionList)),
	)

	respondWithJSON(w, http.StatusOK, dbFactionList)
}

func GetFactionByID(s *state.State, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	factionID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse uuid", err)
		return
	}

	faction, err := services.GetFactionByID(s, r.Context(), factionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestInfo(s, r, "Faction not found in DB")

			respondWithError(w, http.StatusNotFound, "Faction not found in DB", err)
		} else {
			logRequestError(s, r, "failed to fetch faction", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch faction", err)
		}

		return
	}

	logRequestInfo(s, r, "Fetched faction successfully")

	respondWithJSON(w, http.StatusOK, faction)
}
