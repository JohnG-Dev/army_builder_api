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

type UnitHandlers struct {
	S *state.State
}

func (u *UnitHandlers) GetUnits(w http.ResponseWriter, r *http.Request) {

	factionIDString := r.URL.Query().Get("faction_id")

	if factionIDString == "" {
		dbUnitList, err := services.GetUnits(u.S, r.Context(), nil)
		if err != nil {
			logRequestError(u.S, r, "failed to fetch units", err)
			respondWithError(w, http.StatusInternalServerError, "failed to fetch units", err)
			return
		}

		respondWithJSON(w, http.StatusOK, dbUnitList)
		return
	}
	factionID, err := uuid.Parse(factionIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse faction ID", err)
		return
	}
	dbUnitList, err := services.GetUnits(u.S, r.Context(), &factionID)
	if err != nil {
		logRequestError(u.S, r, "failed to fetch Unit List", err)

		respondWithError(w, http.StatusInternalServerError, "failed to fetch Unit List", err)
	}

	logRequestInfo(u.S, r, "Sucessfully fetched factions",
		zap.Int("count", len(dbUnitList)),
	)
	respondWithJSON(w, http.StatusOK, dbUnitList)
}
