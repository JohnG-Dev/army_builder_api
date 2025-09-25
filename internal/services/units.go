package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
	"github.com/google/uuid"
)

func GetUnits(s *state.State, ctx context.Context, factionID uuid.UUID) ([]database.Unit, error) {

	units, err := s.DB.GetUnits(ctx, *factionID)
	if err != nil {
		return nil, err
	}

}

func GetUnitByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Unit, error) {

}
