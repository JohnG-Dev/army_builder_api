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

func mapDBFactionToModel(f database.Faction) models.Faction {
	return models.Faction{
		ID:         f.ID,
		GameID:     f.GameID,
		Name:       f.Name,
		Allegiance: f.Allegiance,
		Version:    f.Version,
		Source:     f.Source,
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
	}
}

func GetFactions(s *state.State, ctx context.Context, gameID *uuid.UUID) ([]models.Faction, error) {
	var dbFactions []database.Faction
	var err error

	if gameID == nil {
		dbFactions, err = s.DB.GetAllFactions(ctx)
	} else {
		dbFactions, err = s.DB.GetFactionsByID(ctx, *gameID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Faction{}, nil
		}
		return nil, err
	}

	if dbFactions == nil {
		return []models.Faction{}, nil
	}

	factions := make([]models.Faction, len(dbFactions))
	for i, f := range dbFactions {
		factions[i] = mapDBFactionToModel(f)
	}

	return factions, nil
}

func GetFactionsByName(s *state.State, ctx context.Context, name string) ([]models.Faction, error) {
	if name == "" {
		return nil, appErr.ErrMissingID
	}

	dbFactions, err := s.DB.GetFactionsByName(ctx, "%"+name+"%")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Faction{}, nil
		}

		return nil, err
	}

	factions := make([]models.Faction, len(dbFactions))
	for i, f := range dbFactions {
		factions[i] = mapDBFactionToModel(f)
	}

	return factions, nil
}

func GetFactionByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Faction, error) {
	if id == uuid.Nil {
		return models.Faction{}, appErr.ErrMissingID
	}

	dbFaction, err := s.DB.GetFaction(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Faction{}, appErr.ErrNotFound
		}
		return models.Faction{}, err
	}

	faction := mapDBFactionToModel(dbFaction)

	return faction, nil
}
