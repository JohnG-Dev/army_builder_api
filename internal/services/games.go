package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetGames(s *state.State, ctx context.Context) ([]models.Game, error) {

	dbGames, err := s.DB.GetGames(ctx)
	if err != nil {
		return nil, err
	}

	if dbGames == nil {
		return []models.Game{}, nil
	}

	games := make([]models.Game, len(dbGames))
	for i, g := range dbGames {
		games[i] = models.Game{
			ID:        g.ID,
			Name:      g.Name,
			Edition:   g.Edition,
			Version:   g.Version,
			Source:    g.Source,
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		}
	}

	return games, nil
}
