package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/JohnG-Dev/army_builder_api/internal/services"
	"github.com/JohnG-Dev/army_builder_api/internal/state"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WeaponsHandlers struct {
	S *state.State
}

func (wh *WeaponsHandlers) GetWeaponsForUnit(w http.ResponseWriter, r *http.Request) {

	unitIDString := r.URL.Query().Get("unit_id")

	if unitIDString == "" {
		respondWithError(w, http.StatusBadRequest, "missing unit id", nil)
		return
	}

	unitID, err := uuid.Parse(unitIDString)
	if err != nil {
		logRequestError(wh.S, r, "failed to parse unit id", err)
		respondWithError(w, http.StatusBadRequest, "failed to parse unit id", err)
		return
	}

	dbWeapons, err := services.GetWeaponsForUnit(wh.S, r.Context(), &unitID)
	if err != nil {
		logRequestError(wh.S, r, "failed to fetch weapons", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch weapons", err)
		return
	}

	logRequestInfo(wh.S, r, "Successfully fetched weapons",
		zap.Int("count", len(dbWeapons)),
	)
	respondWithJSON(w, http.StatusOK, dbWeapons)

}

func (wh *WeaponsHandlers) GetWeaponsByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	weaponID, err := uuid.Parse(id)
	if err != nil {
		logRequestError(wh.S, r, "invalid unit id", err)
		respondWithError(w, http.StatusBadRequest, "invalid unit id", err)
		return
	}

	weapon, err := services.GetWeaponByID(wh.S, r.Context(), weaponID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRequestError(wh.S, r, "weapon not found", err)
			respondWithError(w, http.StatusNotFound, "weapon not found", err)
		} else {
			logRequestError(wh.S, r, "failed to fetch weapon", err)
			respondWithError(w, http.StatusInternalServerError, "failed to fetch weapon", err)
		}
		return
	}

	logRequestInfo(wh.S, r, "Successfully fetched weapon")
	respondWithJSON(w, http.StatusOK, weapon)

}
