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
		weapons[i] = models.Weapon{
			ID:        w.ID,
			UnitID:    w.UnitID,
			Name:      w.Name,
			Range:     w.Range,
			Attacks:   w.Attacks,
			ToHit:     w.ToHit,
			ToWound:   w.ToWound,
			Rend:      w.Rend,
			Damage:    w.Damage,
			Version:   w.Version,
			Source:    w.Source,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		}
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
		weapons[i] = models.Weapon{
			ID:        w.ID,
			UnitID:    w.UnitID,
			Name:      w.Name,
			Range:     w.Range,
			Attacks:   w.Attacks,
			ToHit:     w.ToHit,
			ToWound:   w.ToWound,
			Rend:      w.Rend,
			Damage:    w.Damage,
			Version:   w.Version,
			Source:    w.Source,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		}
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

	weapon := models.Weapon{
		ID:        dbWeapon.ID,
		UnitID:    dbWeapon.UnitID,
		Name:      dbWeapon.Name,
		Range:     dbWeapon.Range,
		Attacks:   dbWeapon.Attacks,
		ToHit:     dbWeapon.ToHit,
		ToWound:   dbWeapon.ToWound,
		Rend:      dbWeapon.Rend,
		Damage:    dbWeapon.Damage,
		Version:   dbWeapon.Version,
		Source:    dbWeapon.Source,
		CreatedAt: dbWeapon.CreatedAt,
		UpdatedAt: dbWeapon.UpdatedAt,
	}

	return weapon, nil
}
