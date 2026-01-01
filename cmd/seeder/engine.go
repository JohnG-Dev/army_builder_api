package main

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type Seeder struct {
	s          *state.State
	ctx        context.Context
	keywordMap map[string]uuid.UUID
	gameMap    map[string]uuid.UUID
}

func NewSeeder(ctx context.Context, s *state.State) *Seeder {
	return &Seeder{
		s:          s,
		ctx:        ctx,
		keywordMap: make(map[string]uuid.UUID),
		gameMap:    make(map[string]uuid.UUID),
	}
}

func (sr *Seeder) SeedFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var seed models.SeedData

	err = yaml.Unmarshal(data, &seed)
	if err != nil {
		return err
	}

	sr.s.Logger.Info("Processing files",
		zap.String("path", path),
		zap.String("game", seed.GameName),
	)

	return nil
}
