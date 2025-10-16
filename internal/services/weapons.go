package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetWeaponsForUnit(s *state.State, ctx context.Context, unitID *uuid.UUID) ([]database.Weapon, error) {

	if unitID == nil {
		return nil, appErr.ErrMissingUnitID
	}

	weapons, err := s.DB.GetWeaponsForUnit(ctx, *unitID)
	if err != nil {
		return nil, err
	}
	if weapons == nil {
		weapons = []database.Weapon{}
	}

	return weapons, nil

}

func GetWeaponByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Weapon, error) {

	if id == uuid.Nil {
		return database.Weapon{}, appErr.ErrMissingID
	}

	weapon, err := s.DB.GetWeaponByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Weapon{}, appErr.ErrNotFound
		}
		return database.Weapon{}, appErr.ErrNotFound
	}

	return weapon, nil
}
