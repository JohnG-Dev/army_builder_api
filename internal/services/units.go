package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetUnits(s *state.State, ctx context.Context, factionID *uuid.UUID) ([]models.Unit, error) {
	var dbUnits []database.Unit
	var err error

	if factionID == nil {
		dbUnits, err = s.DB.GetAllUnits(ctx)
	} else {
		dbUnits, err = s.DB.GetUnits(ctx, *factionID)
	}

	if err != nil {
		return nil, err
	}

	if dbUnits == nil {
		dbUnits = []database.Unit{}
	}

	return dbUnits, nil
}

func GetUnitByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Unit, error) {

	if id == uuid.Nil {
		return database.Unit{}, appErr.ErrMissingID
	}

	unit, err := s.DB.GetUnitByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Unit{}, appErr.ErrNotFound
		}
		return database.Unit{}, err
	}

	return unit, nil
}

func GetManifestations(s *state.State, ctx context.Context) ([]database.Unit, error) {

	manifestations, err := s.DB.GetManifestations(ctx)
	if err != nil {
		return nil, err
	}
	if manifestations == nil {
		manifestations = []database.Unit{}
	}
	return manifestations, nil
}

func GetNonManifestationUnits(s *state.State, ctx context.Context) ([]database.Unit, error) {

	units, err := s.DB.GetNonManifestationUnits(ctx)
	if err != nil {
		return nil, err
	}
	if units == nil {
		units = []database.Unit{}
	}

	return units, nil
}

func GetManifestationByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Unit, error) {

	if id == uuid.Nil {
		return database.Unit{}, appErr.ErrMissingID
	}

	manifestation, err := s.DB.GetManifestationByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Unit{}, appErr.ErrNotFound
		}
		return database.Unit{}, err
	}

	return manifestation, nil
}
