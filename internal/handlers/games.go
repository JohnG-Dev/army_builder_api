package handlers

import (
	"net/http"

	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"

	"go.uber.org/zap"
)

func GetGames(s *state.State, w http.ResponseWriter, r *http.Request) {
	dbGames, err := services.GetGames(s, r.Context())
	if err != nil {

		LogError(s, r, "Failed to fetch games", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't access database", err)
		return
	}

	LogInfo(s, r, "Fetched games successfully",
		zap.Int("count", len(dbGames)),
	)
	respondWithJSON(w, http.StatusOK, dbGames)
}
