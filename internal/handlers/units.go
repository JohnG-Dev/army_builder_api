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

type UnitsHandlers struct {
	S *state.State
}

func (h *UnitsHandlers) GetUnits(w http.ResponseWriter, r *http.Request) {
	factionIDStr := r.URL.Query().Get("faction_id")
	var units []models.Unit
	var err error

	if factionIDStr == "" {
		units, err = services.GetUnits(h.S, r.Context(), nil)
	} else {
		factionID, parseErr := uuid.Parse(factionIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid faction id", parseErr)
			return
		}
		units, err = services.GetUnits(h.S, r.Context(), &factionID)
	}

	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "units not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch units", err)
		}
		logRequestError(h.S, r, "failed to fetch units", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched units", zap.Int("count", len(units)))
	respondWithJSON(w, http.StatusOK, units)
}

func (h *UnitsHandlers) GetUnitByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
		return
	}

	unit, err := services.GetUnitByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "unit not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch unit", err)
		}

		logRequestError(h.S, r, "failed to fetch unit", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched unit")
	respondWithJSON(w, http.StatusOK, unit)
}

func (h *UnitsHandlers) GetManifestations(w http.ResponseWriter, r *http.Request) {

	manifestations, err := services.GetManifestations(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "manifestations not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch manifestations", err)
		}

		logRequestError(h.S, r, "failed to fetch manifestations", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched manifestations", zap.Int("count", len(manifestations)))
	respondWithJSON(w, http.StatusOK, manifestations)
}

func (h *UnitsHandlers) GetNonManifestationUnits(w http.ResponseWriter, r *http.Request) {

	units, err := services.GetNonManifestationUnits(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "units not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch non-manifestation units", err)
		}

		logRequestError(h.S, r, "failed to fetch non-manifestation units", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched non-manifestation units", zap.Int("count", len(units)))
	respondWithJSON(w, http.StatusOK, units)
}

func (h *UnitsHandlers) GetManifestationByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid manifestation id", err)
		return
	}

	manifestation, err := services.GetManifestationByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "manifestation not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch manifestation", err)
		}
		logRequestError(h.S, r, "failed to fetch manifestation", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched manifestation")
	respondWithJSON(w, http.StatusOK, manifestation)
}
