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

type RulesHandlers struct {
	S *state.State
}

func (h *RulesHandlers) GetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := services.GetAllRules(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "rules not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch rules", err)
		}

		logRequestError(h.S, r, "failed to fetch rules", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched rules", zap.Int("count", len(rules)))
	respondWithJSON(w, http.StatusOK, rules)
}

func (h *RulesHandlers) GetRulesForGame(w http.ResponseWriter, r *http.Request) {
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

	rules, err := services.GetRulesForGame(h.S, r.Context(), gameID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing game id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "rules not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch rules", err)

		}

		logRequestError(h.S, r, "failed to fetch rules", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched rules", zap.Int("count", len(rules)))
	respondWithJSON(w, http.StatusOK, rules)
}

func (h *RulesHandlers) GetRulesByType(w http.ResponseWriter, r *http.Request) {
	gameStr := r.URL.Query().Get("game_id")
	typeStr := r.URL.Query().Get("type")

	if gameStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing game id", nil)
		return
	}

	if typeStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing type paramter", nil)
		return
	}

	gameID, err := uuid.Parse(gameStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid game id", err)
	}

	rules, err := services.GetRulesByType(h.S, r.Context(), gameID, typeStr)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "rules not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch rules", err)
		}

		logRequestError(h.S, r, "failed to fetch rules", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched rules", zap.Int("count", len(rules)))
	respondWithJSON(w, http.StatusOK, rules)
}

func (h *RulesHandlers) GetRuleByID(w http.ResponseWriter, r *http.Request) {
	ruleStr := r.PathValue("id")

	if ruleStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing rule id", nil)
		return
	}

	ruleID, err := uuid.Parse(ruleStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid rule id", err)
		return
	}

	rule, err := services.GetRuleByID(h.S, r.Context(), ruleID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "rule not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch rules", err)
		}

		logRequestError(h.S, r, "failed to fetch rule", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched rule")
	respondWithJSON(w, http.StatusOK, rule)
}
