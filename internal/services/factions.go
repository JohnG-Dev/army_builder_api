package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
	"github.com/google/uuid"
)

func ListFactions(s *state.State, ctx context.Context, gameID *uuid.UUID) ([]database.Faction, error) {

	if gameID == nil {
		factions, err := s.DB.GetAllFactions(ctx)

		if err != nil {
			return nil, err
		}

		if factions == nil {
			return []database.Faction{}, nil
		}

		return factions, nil
	}

	factions, err := s.DB.GetFactionsByID(ctx, *gameID)
	if err != nil {
		return nil, err
	}

	if factions == nil {
		return []database.Faction{}, nil
	}

	return factions, nil

}

func GetFactionByID(s *state.State, ctx context.Context, id uuid.UUID) (database.Faction, error) {

	faction, err := s.DB.GetFaction(ctx, id)
	if err != nil {
		return database.Faction{}, err
	}

	return faction, nil
}
