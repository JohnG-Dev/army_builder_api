package handlers

import (
	"net/http"

	"github.com/JohnG-Dev/army_builder_api/internal/api"
)

func (cfg *api.APIConfig) GetGames(w http.ResponseWriter, r *http.Request) {
	dbGames, err := services.GetGames()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't access database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbGames)
}
