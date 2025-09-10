package handlers

import (
	"net/http"

	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"

	"go.uber.org/zap"
)

type GamesHandlers struct {
	S *state.State
}

func (h *GamesHandlers) GetGames(w http.ResponseWriter, r *http.Request) {
	dbGames, err := services.GetGames(h.S, r.Context())
	if err != nil {

		logRequestError(h.S, r, "Failed to fetch games", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	logRequestInfo(h.S, r, "Fetched games successfully",
		zap.Int("count", len(dbGames)),
	)
	respondWithJSON(w, http.StatusOK, dbGames)
}
