package services

import (
	"context"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetGames(s *state.State, ctx context.Context) ([]database.Game, error) {
	return s.DB.GetGames(ctx)
}
