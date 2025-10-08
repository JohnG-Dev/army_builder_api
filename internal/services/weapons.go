package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/state"

	"github.com/google/uuid"
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
		return []database.Weapon{}, nil
	}

	return weapons, nil

}

func GetWeaponByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Weapon, error) {

	weapon, err := s.DB.GetWeaponByID(ctx, id)
	if err != nil {
		return database.Weapon{}, err
	}

	return weapon, nil
}
