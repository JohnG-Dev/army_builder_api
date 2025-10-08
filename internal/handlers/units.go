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

type UnitsHandlers struct {
	S *state.State
}

func (u *UnitsHandlers) GetUnits(w http.ResponseWriter, r *http.Request) {

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
		logRequestError(u.S, r, "failed to fetch units", err)

		respondWithError(w, http.StatusInternalServerError, "failed to fetch units", err)
		return
	}

	logRequestInfo(u.S, r, "Successfully fetched units",
		zap.Int("count", len(dbUnitList)),
	)
	respondWithJSON(w, http.StatusOK, dbUnitList)
}

func (u *UnitsHandlers) GetUnitByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	unitID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
		return
	}

	unit, err := services.GetUnitByID(u.S, r.Context(), unitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestError(u.S, r, "unit not found", err)

			respondWithError(w, http.StatusNotFound, "unit not found", err)
		} else {
			logRequestError(u.S, r, "failed to fetch unit", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch unit", err)
		}

		return
	}

	logRequestInfo(u.S, r, "Successfully fetched unit")

	respondWithJSON(w, http.StatusOK, unit)
}
