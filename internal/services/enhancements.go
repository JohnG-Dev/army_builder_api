package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetEnhancements(s *state.State, ctx context.Context) ([]models.Enhancement, error) {
	dbEnhancements, err := s.DB.GetEnhancements(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Enhancement{}, nil
		}
		return nil, err
	}
	if dbEnhancements == nil {
		return []models.Enhancement{}, nil
	}

	enhancements := make([]models.Enhancement, len(dbEnhancements))
	for i, e := range dbEnhancements {
		enhancements[i] = models.Enhancement{
			ID:              e.ID,
			FactionID:       e.FactionID,
			Name:            e.Name,
			EnhancementType: e.EnhancementType,
			Description:     e.Description,
			Points:          int(e.Points),
			Version:         e.Version,
			Source:          e.Source,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		}
	}

	return enhancements, nil
}

func GetEnhancementsByFaction(s *state.State, ctx context.Context, factionID *uuid.UUID) ([]models.Enhancement, error) {
	if factionID == nil {
		return nil, appErr.ErrMissingFactionID
	}

	dbEnhancements, err := s.DB.GetEnhancementsForFaction(ctx, *factionID)
	if err != nil {
		return nil, err
	}

	if dbEnhancements == nil {
		return []models.Enhancement{}, nil
	}

	enhancements := make([]models.Enhancement, len(dbEnhancements))
	for i, e := range dbEnhancements {
		enhancements[i] = models.Enhancement{
			ID:              e.ID,
			FactionID:       e.FactionID,
			Name:            e.Name,
			EnhancementType: e.EnhancementType,
			Description:     e.Description,
			Points:          int(e.Points),
			Version:         e.Version,
			Source:          e.Source,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		}
	}

	return enhancements, nil
}

func GetEnhancementByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Enhancement, error) {
	if id == uuid.Nil {
		return models.Enhancement{}, appErr.ErrMissingID
	}

	dbEnhancement, err := s.DB.GetEnhancementByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Enhancement{}, appErr.ErrNotFound
		}
		return models.Enhancement{}, err
	}

	enhancement := models.Enhancement{
		ID:              dbEnhancement.ID,
		FactionID:       dbEnhancement.FactionID,
		Name:            dbEnhancement.Name,
		EnhancementType: dbEnhancement.EnhancementType,
		Description:     dbEnhancement.Description,
		Points:          int(dbEnhancement.Points),
		Version:         dbEnhancement.Version,
		Source:          dbEnhancement.Source,
		CreatedAt:       dbEnhancement.CreatedAt,
		UpdatedAt:       dbEnhancement.UpdatedAt,
	}

	return enhancement, nil
}

func GetEnhancementsByType(s *state.State, ctx context.Context, enhancementType string) ([]models.Enhancement, error) {
	if enhancementType == "" {
		return nil, appErr.ErrMissingID
	}

	dbEnhancements, err := s.DB.GetEnhancementsByType(ctx, enhancementType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Enhancement{}, nil
		}

		return nil, err
	}

	if dbEnhancements == nil {
		return []models.Enhancement{}, nil
	}

	enhancements := make([]models.Enhancement, len(dbEnhancements))
	for i, e := range dbEnhancements {
		enhancements[i] = models.Enhancement{
			ID:              e.ID,
			FactionID:       e.FactionID,
			Name:            e.Name,
			EnhancementType: e.EnhancementType,
			Description:     e.Description,
			Points:          int(e.Points),
			Version:         e.Version,
			Source:          e.Source,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		}
	}

	return enhancements, nil
}
