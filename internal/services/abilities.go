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

func GetAbilityEffectsForAbility(s *state.State, ctx context.Context, abilityID uuid.UUID) ([]models.AbilityEffect, error) {

	if abilityID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	dbEffects, err := s.DB.GetAbilityEffectsForAbility(ctx, abilityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.AbilityEffect{}, nil
		}

		return nil, err
	}

	if dbEffects == nil {
		return []models.AbilityEffect{}, nil
	}

	effects := make([]models.AbilityEffect, len(dbEffects))
	for i, e := range dbEffects {
		effects[i] = models.AbilityEffect{
			ID:          e.ID,
			AbilityID:   e.AbilityID,
			Stat:        e.Stat,
			Modifier:    int(e.Modifier),
			Condition:   e.Condition,
			Description: e.Description,
			Version:     e.Version,
			Source:      e.Source,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		}
	}

	return effects, nil
}

func GetAbilitiesForUnit(s *state.State, ctx context.Context, unitID uuid.UUID) ([]models.Ability, error) {

	if unitID == uuid.Nil {
		return nil, appErr.ErrMissingUnitID
	}

	dbAbilities, err := s.DB.GetAbilitiesForUnit(ctx, unitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Ability{}, nil
		}
		return nil, err
	}

	if dbAbilities == nil {
		return []models.Ability{}, nil
	}

	abilities := make([]models.Ability, len(dbAbilities))
	for i, a := range dbAbilities {
		effects, _ := GetAbilityEffectsForAbility(s, ctx, a.ID)

		abilities[i] = models.Ability{
			ID:          a.ID,
			UnitID:      a.UnitID,
			FactionID:   a.FactionID,
			Name:        a.Name,
			Type:        a.Type,
			Phase:       a.Phase,
			Description: a.Description,
			Version:     a.Version,
			Source:      a.Source,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
			Effects:     effects,
		}
	}

	return abilities, nil
}

func GetAbilitiesForFaction(s *state.State, ctx context.Context, factionID uuid.UUID) ([]models.Ability, error) {

	if factionID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	dbAbilities, err := s.DB.GetAbilitiesForFaction(ctx, factionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Ability{}, nil
		}

		return nil, err
	}

	if dbAbilities == nil {
		return []models.Ability{}, nil
	}

	abilities := make([]models.Ability, len(dbAbilities))
	for i, a := range dbAbilities {
		effects, _ := GetAbilityEffectsForAbility(s, ctx, a.ID)

		abilities[i] = models.Ability{
			ID:          a.ID,
			UnitID:      a.UnitID,
			FactionID:   a.FactionID,
			Name:        a.Name,
			Type:        a.Type,
			Phase:       a.Phase,
			Description: a.Description,
			Version:     a.Version,
			Source:      a.Source,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
			Effects:     effects,
		}
	}

	return abilities, nil
}

func GetAbilityByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Ability, error) {

	if id == uuid.Nil {
		return models.Ability{}, appErr.ErrMissingID
	}

	dbAbility, err := s.DB.GetAbilityByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Ability{}, appErr.ErrNotFound
		}
		return models.Ability{}, err
	}

	effects, _ := GetAbilityEffectsForAbility(s, ctx, dbAbility.ID)

	ability := models.Ability{
		ID:          dbAbility.ID,
		UnitID:      dbAbility.UnitID,
		FactionID:   dbAbility.FactionID,
		Name:        dbAbility.Name,
		Type:        dbAbility.Type,
		Phase:       dbAbility.Phase,
		Description: dbAbility.Description,
		Version:     dbAbility.Version,
		Source:      dbAbility.Source,
		CreatedAt:   dbAbility.CreatedAt,
		UpdatedAt:   dbAbility.UpdatedAt,
		Effects:     effects,
	}

	return ability, nil
}
