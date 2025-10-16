package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetFactions(s *state.State, ctx context.Context, gameID *uuid.UUID) ([]database.Faction, error) {
	var factions []database.Faction
	var err error

	if gameID == nil {
		factions, err = s.DB.GetAllFactions(ctx)
	} else {
		factions, err = s.DB.GetFactionsByID(ctx, *gameID)
	}

	if err != nil {
		return nil, err
	}

	if factions == nil {
		factions = []database.Faction{}
	}

	return factions, nil
}

func GetFactionByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Faction, error) {

	if id == uuid.Nil {
		return database.Faction{}, appErr.ErrMissingID
	}

	faction, err := s.DB.GetFaction(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Faction{}, appErr.ErrNotFound
		}
		return database.Faction{}, err
	}

	return faction, nil
}
