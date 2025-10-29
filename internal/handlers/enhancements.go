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

type EnhancementsHandlers struct {
	S *state.State
}

func (h *EnhancementsHandlers) GetEnhancements(w http.ResponseWriter, r *http.Request) {
	factionIDStr := r.URL.Query().Get("faction_id")
	var enhancements []models.Enhancement
	var err error

	if factionIDStr == "" {
		enhancements, err = services.GetEnhancements(h.S, r.Context())
	} else {
		factionID, parseErr := uuid.Parse(factionIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid faction id", parseErr)
			return
		}
		enhancements, err = services.GetEnhancementsByFaction(h.S, r.Context(), &factionID)
	}

	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingFactionID):
			respondWithError(w, http.StatusBadRequest, "missing faction id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "enhancements not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch enhancements", err)
		}

		logRequestError(h.S, r, "failed to fetch enhancements", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched enhancements", zap.Int("count", len(enhancements)))
	respondWithJSON(w, http.StatusOK, enhancements)
}

func (h *EnhancementsHandlers) GetEnhancementByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid enhancement id", err)
		return
	}

	enhancement, err := services.GetEnhancementByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "enhancement not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch enhancement", err)
		}

		logRequestError(h.S, r, "failed to fetch enhancement", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched enhancement")
	respondWithJSON(w, http.StatusOK, enhancement)
}

func (h *EnhancementsHandlers) GetEnhancementsByType(w http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")

	if typeStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing enhancements type", nil)
		return
	}

	enhancements, err := services.GetEnhancementsByType(h.S, r.Context(), typeStr)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing enhancement type", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "enhancements not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch enhancements", err)
		}

		logRequestError(h.S, r, "failed to fetch enhancements", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched ehnahncements", zap.Int("count", len(enhancements)))
	respondWithJSON(w, http.StatusOK, enhancements)
}
