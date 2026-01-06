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

func mapDBGameToModel(g database.Game) models.Game {
	return models.Game{
		ID:        g.ID,
		Name:      g.Name,
		Edition:   g.Edition,
		Version:   g.Version,
		Source:    g.Source,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

func GetGames(s *state.State, ctx context.Context) ([]models.Game, error) {
	dbGames, err := s.DB.GetGames(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Game{}, nil
		}
		return nil, err
	}

	if dbGames == nil {
		return []models.Game{}, nil
	}

	games := make([]models.Game, len(dbGames))
	for i, g := range dbGames {
		games[i] = mapDBGameToModel(g)
	}

	return games, nil
}

func GetGame(s *state.State, ctx context.Context, id uuid.UUID) (models.Game, error) {
	if id == uuid.Nil {
		return models.Game{}, appErr.ErrMissingID
	}

	dbGame, err := s.DB.GetGame(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Game{}, appErr.ErrNotFound
		}

		return models.Game{}, err

	}
	game := mapDBGameToModel(dbGame)

	return game, nil
}

func GetGameByName(s *state.State, ctx context.Context, name string) (models.Game, error) {
	if name == "" {
		return models.Game{}, appErr.ErrMissingID
	}

	dbGame, err := s.DB.GetGameByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Game{}, appErr.ErrNotFound
		}

		return models.Game{}, err
	}

	game := mapDBGameToModel(dbGame)
	return game, nil
}
