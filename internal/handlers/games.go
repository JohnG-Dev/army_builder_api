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

type GamesHandlers struct {
	S *state.State
}

func (h *GamesHandlers) GetGames(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name != "" {
		h.getGameByName(w, r)
		return
	}

	h.getAllGames(w, r)
}

func (h *GamesHandlers) getAllGames(w http.ResponseWriter, r *http.Request) {
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

func (h *GamesHandlers) getGameByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		respondWithError(w, http.StatusBadRequest, "missing name", nil)
		return
	}

	game, err := services.GetGameByName(h.S, r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "game not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch game", err)
		}

		logRequestError(h.S, r, "failed to fetch game", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched game")
	respondWithJSON(w, http.StatusOK, game)
}

func (h *GamesHandlers) GetGameByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing id", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", err)
		return
	}

	game, err := services.GetGame(h.S, r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrNotFound):
			respondWithError(w, http.StatusNotFound, "game not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to fetch game", err)
		}

		logRequestError(h.S, r, "failed to fetch game", err)
		return
	}

	logRequestInfo(h.S, r, "Successfully fetched game")
	respondWithJSON(w, http.StatusOK, game)
}
