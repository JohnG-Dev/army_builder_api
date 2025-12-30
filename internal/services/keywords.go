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

func GetAllKeywords(s *state.State, ctx context.Context) ([]models.Keyword, error) {
	dbKeywords, err := s.DB.GetAllKeywords(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Keyword{}, nil
		}

		return nil, err
	}

	keywords := make([]models.Keyword, len(dbKeywords))
	for i, k := range dbKeywords {
		keywords[i] = models.Keyword{
			ID:          k.ID,
			GameID:      k.GameID,
			Name:        k.Name,
			Description: k.Description,
			Version:     k.Version,
			Source:      k.Source,
			CreatedAt:   k.CreatedAt,
			UpdatedAt:   k.UpdatedAt,
		}
	}

	return keywords, nil
}

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
			UnitID:      uk.UnitID,
			KeywordID:   uk.KeywordID,
			KeywordName: uk.KeywordName,
			Value:       uk.Value,
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
			MinUnitSize:     int(u.MinUnitSize),
			MaxUnitSize:     int(u.MaxUnitSize),
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

	dbUnits, err := s.DB.GetUnitsWithKeywordAndValue(ctx, database.GetUnitsWithKeywordAndValueParams{
		Name:  name,
		Value: value,
	})
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
			MinUnitSize:     int(u.MinUnitSize),
			MaxUnitSize:     int(u.MaxUnitSize),
			MatchedPlay:     u.MatchedPlay,
			Version:         u.Version,
			Source:          u.Source,
			CreatedAt:       u.CreatedAt,
			UpdatedAt:       u.UpdatedAt,
		}
	}

	return units, nil
}

func GetKeywordsForGame(s *state.State, ctx context.Context, gameID uuid.UUID) ([]models.Keyword, error) {
	if gameID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	dbKeywords, err := s.DB.GetKeywordsForGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Keyword{}, nil
		}

		return nil, err
	}

	if dbKeywords == nil {
		return []models.Keyword{}, nil
	}

	keywords := make([]models.Keyword, len(dbKeywords))
	for i, k := range dbKeywords {
		keywords[i] = models.Keyword{
			ID:          k.ID,
			GameID:      k.GameID,
			Name:        k.Name,
			Description: k.Description,
			Version:     k.Version,
			Source:      k.Source,
			CreatedAt:   k.CreatedAt,
			UpdatedAt:   k.UpdatedAt,
		}
	}

	return keywords, nil
}

func GetKeywordByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Keyword, error) {
	if id == uuid.Nil {
		return models.Keyword{}, appErr.ErrMissingID
	}

	dbKeyword, err := s.DB.GetKeywordByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Keyword{}, appErr.ErrNotFound
		}

		return models.Keyword{}, err
	}

	keyword := models.Keyword{
		ID:          dbKeyword.ID,
		GameID:      dbKeyword.GameID,
		Name:        dbKeyword.Name,
		Description: dbKeyword.Description,
		Version:     dbKeyword.Version,
		Source:      dbKeyword.Source,
		CreatedAt:   dbKeyword.CreatedAt,
		UpdatedAt:   dbKeyword.UpdatedAt,
	}

	return keyword, nil
}
