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

		s.Logger.Error("Failed to fetch games",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)
		respondWithError(w, http.StatusInternalServerError, "Couldn't access database", err)
		return
	}

	s.Logger.Info("Fetched games successfully",
		zap.Int("count", len(dbGames)),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)
	respondWithJSON(w, http.StatusOK, dbGames)
}
