package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type GamesHandlers struct {
	S *state.State
}

func (h *GamesHandlers) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := services.GetGames(h.S, r.Context())
	if err != nil {

		logRequestError(h.S, r, "failed to fetch games", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched games",
		zap.Int("count", len(games)),
	)
	respondWithJSON(w, http.StatusOK, games)
}
