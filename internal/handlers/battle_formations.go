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

type BattleFormationsHandlers struct {
	S *state.State
}

func (h *BattleFormationsHandlers) GetBattleFormations(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	factionID := r.URL.Query().Get("faction_id")

	if gameID != "" {
		h.getBattleFormationsForGame(w, r)
		return
	}

	if factionID != "" {
		h.getBattleFormationsForFaction(w, r)
		return
	}

	h.getAllBattleFormations(w, r)
}

func (h *BattleFormationsHandlers) getAllBattleFormations(w http.ResponseWriter, r *http.Request) {
	battleFormations, err := services.GetAllBattleFormations(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "unable to find battle formations", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch battle formations", err)
		}

		logRequestError(h.S, r, "failed to fetch battle formations", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched battle formations", zap.Int("count", len(battleFormations)))
	respondWithJSON(w, http.StatusOK, battleFormations)
}

func (h *BattleFormationsHandlers) getBattleFormationsForGame(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.URL.Query().Get("game_id")

	if gameIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing game id", nil)
		return
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid game id", err)
		return
	}

	battleFormations, err := services.GetBattleFormationsForGame(h.S, r.Context(), gameID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingID):
			respondWithError(w, http.StatusBadRequest, "missing game id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "unable to find battle formations", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch battle formations", err)
		}

		logRequestError(h.S, r, "failed to fetch battle formations", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched battle formations", zap.Int("count", len(battleFormations)))
	respondWithJSON(w, http.StatusOK, battleFormations)
}

func (h *BattleFormationsHandlers) getBattleFormationsForFaction(w http.ResponseWriter, r *http.Request) {
	factionIDStr := r.URL.Query().Get("faction_id")

	if factionIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing faction id", nil)
		return
	}

	factionID, err := uuid.Parse(factionIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid faction id", err)
		return
	}

	battleFormations, err := services.GetBattleFormationsForFaction(h.S, r.Context(), factionID)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrMissingFactionID):
			respondWithError(w, http.StatusBadRequest, "missing faction id", err)
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "unable to find battle formations", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch battle formations", err)
		}

		logRequestError(h.S, r, "failed to fetch battle formations", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched battle formations", zap.Int("count", len(battleFormations)))
	respondWithJSON(w, http.StatusOK, battleFormations)
}

func (h *BattleFormationsHandlers) GetBattleFormationByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing battle formation id", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid battle formation id", err)
		return
	}

	battleFormation, err := services.GetBattleFormationByID(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "unable to find battle formation", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch battle formation", err)
		}

		logRequestError(h.S, r, "failed to fetch battle formation", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched battle formation")
	respondWithJSON(w, http.StatusOK, battleFormation)
}
