package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type GamesHandlers struct {
	S *state.State
}

func (h *GamesHandlers) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := services.GetGames(h.S, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "games not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch games", err)
		}

		logRequestError(h.S, r, "failed to fetch games", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched games",
		zap.Int("count", len(games)),
	)
	respondWithJSON(w, http.StatusOK, games)
}
