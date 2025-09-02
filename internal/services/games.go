package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetGames(s *state.State, ctx context.Context) ([]database.Game, error) {
	games, err := s.DB.GetGames(ctx)
	if err != nil {
		return nil, err
	}

	if games == nil {
		return []database.Game{}, nil
	}

	return games, nil
}
