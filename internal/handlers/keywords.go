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

type KeywordsHandlers struct {
	S *state.State
}

func (h *KeywordsHandlers) GetKeywords(w http.ResponseWriter, r *http.Request) {
	keywords, err := services.GetAllKeywords(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "keywords not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch keywords", err)
		}

		logRequestError(h.S, r, "failed to fetch keywords", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched keywords", zap.Int("count", len(keywords)))
	respondWithJSON(w, http.StatusOK, keywords)
}

func (h *KeywordsHandlers) GetKeywordsByGame(w http.ResponseWriter, r *http.Request) {
	gameStr := r.URL.Query().Get("game_id")

	if gameStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing game id", nil)
		return
	}

	gameID, err := uuid.Parse(gameStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid game id", err)
		return
	}

	keywords, err := services.GetKeywordsByGame(h.S, r.Context(), gameID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing game id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "keywords not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch keywords", err)
		}

		logRequestError(h.S, r, "failed to fetch keywords", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched keywords", zap.Int("count", len(keywords)))
	respondWithJSON(w, http.StatusOK, keywords)
}

func (h *KeywordsHandlers) GetKeywordsForUnit(w http.ResponseWriter, r *http.Request) {
	unitStr := r.URL.Query().Get("unit_id")

	if unitStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing unit id", nil)
		return
	}

	unitID, err := uuid.Parse(unitStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
		return
	}

	keywords, err := services.GetKeywordsForUnit(h.S, r.Context(), unitID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing unit id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "keywords not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch keywords", err)
		}

		logRequestError(h.S, r, "failed to fetch keywords", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched keywords", zap.Int("count", len(keywords)))
	respondWithJSON(w, http.StatusOK, keywords)
}

func (h *KeywordsHandlers) GetUnitsWithKeyword(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if name == "" {
		respondWithError(w, http.StatusBadRequest, "missing keyword", nil)
		return
	}

	units, err := services.GetUnitsWithKeyword(h.S, r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "units with keyword not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch units with keyword", err)
		}

		logRequestError(h.S, r, "failed to fetch units with keyword", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched units with keyword", zap.Int("count", len(units)))
	respondWithJSON(w, http.StatusOK, units)
}

func (h *KeywordsHandlers) GetUnitsWithKeywordAndValue(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	value := r.PathValue("value")

	if name == "" {
		respondWithError(w, http.StatusBadRequest, "missing name", nil)
		return
	}

	if value == "" {
		respondWithError(w, http.StatusBadRequest, "missing value", nil)
		return
	}

	units, err := services.GetUnitsWithKeywordAndValue(h.S, r.Context(), name, value)
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

func (h *KeywordsHandlers) GetKeywordByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing keyword id", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid keyword id", err)
		return
	}

	keyword, err := services.GetKeywordByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing keyword id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "keyword not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch keyword", err)
		}

		logRequestError(h.S, r, "failed to fetch keyword", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched keyword")
	respondWithJSON(w, http.StatusOK, keyword)
}
