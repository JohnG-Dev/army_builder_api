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

func mapDBWeaponToModel(w database.Weapon) models.Weapon {
	return models.Weapon{
		ID:            w.ID,
		UnitID:        w.UnitID,
		Name:          w.Name,
		Range:         w.Range,
		Attacks:       w.Attacks,
		HitStats:      w.HitStats,
		WoundStrength: w.WoundStrength,
		RendAP:        w.RendAp,
		Damage:        w.Damage,
		Version:       w.Version,
		Source:        w.Source,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}

func GetAllWeapons(s *state.State, ctx context.Context) ([]models.Weapon, error) {
	dbWeapons, err := s.DB.GetAllWeapons(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Weapon{}, nil
		}

		return nil, err
	}

	if dbWeapons == nil {
		return []models.Weapon{}, nil
	}

	weapons := make([]models.Weapon, len(dbWeapons))
	for i, w := range dbWeapons {
		weapons[i] = mapDBWeaponToModel(w)
	}

	return weapons, nil
}

func GetWeaponsForUnit(s *state.State, ctx context.Context, unitID *uuid.UUID) ([]models.Weapon, error) {
	if unitID == nil {
		return nil, appErr.ErrMissingUnitID
	}

	dbWeapons, err := s.DB.GetWeaponsForUnit(ctx, *unitID)
	if err != nil {
		return nil, err
	}
	if dbWeapons == nil {
		return []models.Weapon{}, nil
	}

	weapons := make([]models.Weapon, len(dbWeapons))
	for i, w := range dbWeapons {
		weapons[i] = mapDBWeaponToModel(w)
	}

	return weapons, nil
}

func GetWeaponByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Weapon, error) {
	if id == uuid.Nil {
		return models.Weapon{}, appErr.ErrMissingID
	}

	dbWeapon, err := s.DB.GetWeaponByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Weapon{}, appErr.ErrNotFound
		}
		return models.Weapon{}, err
	}

	weapon := mapDBWeaponToModel(dbWeapon)

	return weapon, nil
}
