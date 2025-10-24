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

func GetKeywordsForUnit(s *state.State, ctx context.Context, unitID uuid.UUID) ([]models.UnitKeyword, error) {

	if unitID == uuid.Nil {
		return nil, appErr.ErrMissingUnitID
	}

	dbUnitKeywords, err := s.DB.GetKeywordsForUnit(ctx, unitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.UnitKeyword{}, nil
		}

		return nil, err
	}

	if dbUnitKeywords == nil {
		return []models.UnitKeyword{}, nil
	}

	keywords := make([]models.UnitKeyword, len(dbUnitKeywords))
	for i, uk := range dbUnitKeywords {
		keywords[i] = models.UnitKeyword{
			UnitID:    uk.UnitID,
			KeywordID: uk.KeywordID,
			Value:     uk.Value,
		}
	}

	return keywords, nil
}

func GetUnitsWithKeyword(s *state.State, ctx context.Context, name string) ([]models.Unit, error) {

	if name == "" {
		return nil, appErr.ErrMissingID
	}

	dbUnits, err := s.DB.GetUnitsWithKeyword(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Unit{}, nil
		}

		return nil, err
	}

	if dbUnits == nil {
		return []models.Unit{}, nil
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = models.Unit{
			ID:              u.ID,
			FactionID:       u.FactionID,
			Name:            u.Name,
			Description:     u.Description,
			IsManifestation: u.IsManifestation,
			Move:            int(u.Move),
			Health:          int(u.Health),
			Save:            u.Save,
			Ward:            u.Ward,
			Control:         int(u.Control),
			Points:          int(u.Points),
			SummonCost:      u.SummonCost,
			Banishment:      u.Banishment,
			MinSize:         int(u.MinSize),
			MaxSize:         int(u.MaxSize),
			MatchedPlay:     u.MatchedPlay,
			Version:         u.Version,
			Source:          u.Source,
			CreatedAt:       u.CreatedAt,
			UpdatedAt:       u.UpdatedAt,
		}
	}

	return units, nil
}

func GetUnitsWithKeywordAndValue(s *state.State, ctx context.Context, name string, value string) ([]models.Unit, error) {

	if name == "" {
		return nil, appErr.ErrMissingID
	}

	dbUnits, err := s.DB.GetUnitsWithKeywordAndValue(ctx, name, value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Unit{}, nil
		}
		return nil, err
	}

	if dbUnits == nil {
		return []models.Unit{}, nil
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = models.Unit{
			ID:              u.ID,
			FactionID:       u.FactionID,
			Name:            u.Name,
			Description:     u.Description,
			IsManifestation: u.IsManifestation,
			Move:            int(u.Move),
			Health:          int(u.Health),
			Save:            u.Save,
			Ward:            u.Ward,
			Control:         int(u.Control),
			Points:          int(u.Points),
			SummonCost:      u.SummonCost,
			Banishment:      u.Banishment,
			MinSize:         int(u.MinSize),
			MaxSize:         int(u.MaxSize),
			MatchedPlay:     u.MatchedPlay,
			Version:         u.Version,
			Source:          u.Source,
			CreatedAt:       u.CreatedAt,
			UpdatedAt:       u.UpdatedAt,
		}
	}

	return units, nil
}
