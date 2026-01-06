package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	appErr "github.com/JohnG-Dev/army_builder_api/internal/errors"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func mapDBUnitToModel(u database.Unit) models.Unit {
	var addStats map[string]string
	if len(u.AdditionalStats) > 0 {
		_ = json.Unmarshal(u.AdditionalStats, &addStats)
	}
	if addStats == nil {
		addStats = make(map[string]string)
	}

	return models.Unit{
		ID:              u.ID,
		FactionID:       u.FactionID,
		Name:            u.Name,
		IsUnique:        u.IsUnique,
		Description:     u.Description,
		IsManifestation: u.IsManifestation,
		Move:            u.Move,
		HealthWounds:    u.HealthWounds,
		Save:            u.SaveStats,
		WardFNP:         u.WardFnp,
		InvulnSave:      u.InvulnSave,
		ControlOC:       u.ControlOc,
		Toughness:       u.Toughness,
		Leadership:      u.LeadershipBravery,
		AdditionalStats: addStats,
		Points:          int(u.Points),
		SummonCost:      u.SummonCost,
		Banishment:      u.Banishment,
		MinUnitSize:     int(u.MinUnitSize),
		MaxUnitSize:     int(u.MaxUnitSize),
		MatchedPlay:     u.MatchedPlay,
		Version:         u.Version,
		Source:          u.Source,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func GetUnits(s *state.State, ctx context.Context, factionID *uuid.UUID) ([]models.Unit, error) {
	var dbUnits []database.Unit
	var err error

	if factionID == nil {
		dbUnits, err = s.DB.GetAllUnits(ctx)
	} else {
		dbUnits, err = s.DB.GetUnitsByFaction(ctx, *factionID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Unit{}, nil
		}
		return nil, err
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = mapDBUnitToModel(u)
	}

	return units, nil
}

func GetUnitsByFaction(s *state.State, ctx context.Context, factionID uuid.UUID) ([]models.Unit, error) {
	if factionID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}
	return GetUnits(s, ctx, &factionID)
}

func GetUnitByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Unit, error) {
	if id == uuid.Nil {
		return models.Unit{}, appErr.ErrMissingID
	}

	dbUnit, err := s.DB.GetUnitByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Unit{}, appErr.ErrNotFound
		}
		return models.Unit{}, err
	}

	unit := mapDBUnitToModel(dbUnit)

	weapons, _ := GetWeaponsForUnit(s, ctx, &id)
	unit.Weapons = weapons

	abilities, _ := GetAbilitiesForUnit(s, ctx, id)
	unit.Abilities = abilities

	keywords, _ := GetKeywordsForUnit(s, ctx, id)
	unit.Keywords = keywords

	return unit, nil
}

func GetManifestations(s *state.State, ctx context.Context) ([]models.Unit, error) {
	dbManifestations, err := s.DB.GetManifestations(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Unit{}, nil
		}
		return nil, err
	}

	manifestations := make([]models.Unit, len(dbManifestations))
	for i, m := range dbManifestations {
		manifestations[i] = mapDBUnitToModel(m)
	}

	return manifestations, nil
}

func GetNonManifestationUnits(s *state.State, ctx context.Context) ([]models.Unit, error) {
	dbUnits, err := s.DB.GetNonManifestationUnits(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Unit{}, nil
		}
		return nil, err
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = mapDBUnitToModel(u)
	}

	return units, nil
}

func GetManifestationByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Unit, error) {
	if id == uuid.Nil {
		return models.Unit{}, appErr.ErrMissingID
	}

	dbManifestation, err := s.DB.GetManifestationByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Unit{}, appErr.ErrNotFound
		}
		return models.Unit{}, err
	}

	unit := mapDBUnitToModel(dbManifestation)

	weapons, _ := GetWeaponsForUnit(s, ctx, &id)
	unit.Weapons = weapons

	abilities, _ := GetAbilitiesForUnit(s, ctx, id)
	unit.Abilities = abilities

	keywords, _ := GetKeywordsForUnit(s, ctx, id)
	unit.Keywords = keywords

	return unit, nil
}

func GetUnitsByMatchedPlay(s *state.State, ctx context.Context, factionID uuid.UUID) ([]models.Unit, error) {
	if factionID == uuid.Nil {
		return nil, appErr.ErrMissingFactionID
	}

	dbUnits, err := s.DB.GetUnitsByMatchedPlay(ctx, factionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Unit{}, nil
		}
		return nil, err
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = mapDBUnitToModel(u)
	}

	return units, nil
}

