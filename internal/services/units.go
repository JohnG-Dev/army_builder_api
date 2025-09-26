package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
	"github.com/google/uuid"
)

func GetUnits(s *state.State, ctx context.Context, factionID *uuid.UUID) ([]database.Unit, error) {

	if factionID == nil {
		unitsList, err := s.DB.ListUnits(ctx)
		if err != nil {
			return nil, err
		}
		if unitsList == nil {
			return []database.Unit{}, nil
		}
		return unitsList, nil
	}

	units, err := s.DB.GetUnits(ctx, *factionID)
	if err != nil {
		return nil, err
	}
	if units == nil {
		return []database.Unit{}, nil
	}

	return units, nil
}

func GetUnitByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Unit, error) {

	unit, err := s.DB.GetUnitByID(ctx, id)
	if err != nil {
		return database.Unit{}, err
	}

	return unit, nil

}
