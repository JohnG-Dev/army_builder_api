package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func GetAllBattleFormations(s *state.State, ctx context.Context) ([]models.BattleFormation, error) {

	dbBattleFormation, err := s.DB.GetAllBattleFormations(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.BattleFormation{}, nil
		}

		return nil, err
	}

	if dbBattleFormation == nil {
		return []models.BattleFormation{}, nil
	}

	battleFormation := make([]models.BattleFormation, len(dbBattleFormation))
	for i, f := range dbBattleFormation {
		battleFormation[i] = models.BattleFormation{
			ID:          f.ID,
			GameID:      f.GameID,
			FactionID:   f.FactionID,
			Name:        f.Name,
			Description: f.Description,
			Version:     f.Version,
			Source:      f.Source,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		}
	}

	return battleFormation, nil
}

func GetBattleFormationsForGame(s *state.State, ctx context.Context, gameID uuid.UUID) ([]models.BattleFormation, error) {

	if gameID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	dbBattleFormation, err := s.DB.GetBattleFormationsForGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.BattleFormation{}, nil
		}

		return nil, err
	}

	if dbBattleFormation == nil {
		return []models.BattleFormation{}, nil
	}

	battleFormations := make([]models.BattleFormation, len(dbBattleFormation))
	for i, f := range dbBattleFormation {
		battleFormations[i] = models.BattleFormation{
			ID:          f.ID,
			GameID:      f.GameID,
			FactionID:   f.FactionID,
			Name:        f.Name,
			Description: f.Description,
			Version:     f.Version,
			Source:      f.Source,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		}
	}

	return battleFormations, nil
}

func GetBattleFormationsForFaction(s *state.State, ctx context.Context, factionID uuid.UUID) ([]models.BattleFormation, error) {

	if factionID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	dbBattleFormation, err := s.DB.GetBattleFormationsForFaction(ctx, factionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.BattleFormation{}, nil
		}

		return nil, err
	}

	if dbBattleFormation == nil {
		return []models.BattleFormation{}, nil
	}

	battleFormation := make([]models.BattleFormation, len(dbBattleFormation))
	for i, f := range dbBattleFormation {
		battleFormation[i] = models.BattleFormation{
			ID:          f.ID,
			GameID:      f.GameID,
			FactionID:   f.FactionID,
			Name:        f.Name,
			Description: f.Description,
			Version:     f.Version,
			Source:      f.Source,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		}
	}

	return battleFormation, nil
}

func GetBattleFormationByID(s *state.State, ctx context.Context, id uuid.UUID) (models.BattleFormation, error) {

	if id == uuid.Nil {
		return models.BattleFormation{}, appErr.ErrMissingID
	}

	dbBattleFormation, err := s.DB.GetBattleFormationByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BattleFormation{}, nil
		}

		return models.BattleFormation{}, err
	}

	battleFormation := models.BattleFormation{
		ID:          dbBattleFormation.ID,
		GameID:      dbBattleFormation.GameID,
		FactionID:   dbBattleFormation.FactionID,
		Name:        dbBattleFormation.Name,
		Description: dbBattleFormation.Description,
		Version:     dbBattleFormation.Version,
		Source:      dbBattleFormation.Source,
		CreatedAt:   dbBattleFormation.CreatedAt,
		UpdatedAt:   dbBattleFormation.UpdatedAt,
	}

	return battleFormation, nil
}
