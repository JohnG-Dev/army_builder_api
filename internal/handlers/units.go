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
		respondWithError(w, http.StatusBadRequest, "invalid faction id", err)
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

func (u *UnitsHandlers) GetManifestations(w http.ResponseWriter, r *http.Request) {
	dbManifestations, err := services.GetManifestations(u.S, r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestError(u.S, r, "manifestations not found", err)

			respondWithError(w, http.StatusNotFound, "manifestations not found", err)
		} else {
			logRequestError(u.S, r, "failed to fetch manifestations", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch manifestations", err)
		}
		return
	}

	logRequestInfo(u.S, r, "Successfully fetched mainfestations",
		zap.Int("count", len(dbManifestations)),
	)

	respondWithJSON(w, http.StatusOK, dbManifestations)
}

func (u *UnitsHandlers) GetNonManifestationUnits(w http.ResponseWriter, r *http.Request) {
	dbUnits, err := services.GetNonManifestationUnits(u.S, r.Context())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestError(u.S, r, "units not found", err)

			respondWithError(w, http.StatusNotFound, "units not found", err)
		} else {
			logRequestError(u.S, r, "failed to fetch units", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch units", err)
		}
		return
	}

	logRequestInfo(u.S, r, "Successfuly fetched units",
		zap.Int("count", len(dbUnits)),
	)

	respondWithJSON(w, http.StatusOK, dbUnits)
}

func (u *UnitsHandlers) GetManifestationByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	manifestationID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
	}

	manifestation, err := services.GetManifestationByID(u.S, r.Context(), manifestationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestError(u.S, r, "manifestation not found", err)

			respondWithError(w, http.StatusNotFound, "manifestation not found", err)
		} else {
			logRequestError(u.S, r, "failed to fetch manifestation", err)

			respondWithError(w, http.StatusInternalServerError, "failed to fetch manifestation", err)
		}

		return
	}

	logRequestInfo(u.S, r, "Successfuly fetched manifestation")

	respondWithJSON(w, http.StatusOK, manifestation)
}
