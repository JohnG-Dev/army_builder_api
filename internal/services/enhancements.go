package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
	"github.com/google/uuid"
)

func GetEnhancements(s *state.State, ctx context.Context) ([]database.Enhancement, error) {

	enhancements, err := s.DB.GetEnhancements(ctx)
	if err != nil {
		return nil, err
	}
	if enhancements == nil {
		return []database.Enhancement{}, nil
	}

	return enhancements, nil
}

func GetEnhancementsByFaction(s *state.State, ctx context.Context, factionID *uuid.UUID) ([]database.Enhancement, error) {

	if factionID == nil {
		return nil, appErr.ErrMissingFactionID
	}

	enhancements, err := s.DB.GetEnhancementsForFaction(ctx, *factionID)
	if err != nil {
		return nil, err
	}

	if enhancements == nil {
		enhancements = []database.Enhancement{}
	}

	return enhancements, nil
}

func GetEnhancementByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Enhancement, error) {

	if id == uuid.Nil {
		return database.Enhancement{}, appErr.ErrMissingID
	}

	enhancement, err := s.DB.GetEnhancementByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Enhancement{}, appErr.ErrNotFound
		}
		return database.Enhancement{}, err
	}

	return enhancement, nil
}
