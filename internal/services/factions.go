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

func GetFactions(s *state.State, ctx context.Context, gameID *uuid.UUID) ([]models.Faction, error) {
	var dbFactions []database.Faction
	var err error

	if gameID == nil {
		dbFactions, err = s.DB.GetAllFactions(ctx)
	} else {
		dbFactions, err = s.DB.GetFactionsByID(ctx, *gameID)
	}

	if err != nil {
		return nil, err
	}

	if dbFactions == nil {
		return []models.Faction{}, nil
	}

	factions := make([]models.Faction, len(dbFactions))
	for i, f := range dbFactions {
		factions[i] = models.Faction{
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

	faction := models.Faction{
		ID:         dbFaction.ID,
		GameID:     dbFaction.GameID,
		Name:       dbFaction.Name,
		Allegiance: dbFaction.Allegiance,
		Version:    dbFaction.Version,
		Source:     dbFaction.Source,
		CreatedAt:  dbFaction.CreatedAt,
		UpdatedAt:  dbFaction.UpdatedAt,
	}

	return faction, nil
}
