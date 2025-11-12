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

	if dbUnits == nil {
		dbUnits = []database.Unit{}
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = models.Unit{
			ID:              u.ID,
			FactionID:       u.FactionID,
			Name:            u.Name,
			Description:     u.Description,
			IsManifestation: u.IsManifestation,
			Move:            int(u.Move),
			Health:          int(u.Health),
			Save:            u.Save,
			Ward:            u.Ward,
			Control:         int(u.Control),
			Points:          int(u.Points),
			SummonCost:      u.SummonCost,
			Banishment:      u.Banishment,
			MinSize:         int(u.MinSize),
			MaxSize:         int(u.MaxSize),
			MatchedPlay:     u.MatchedPlay,
			Version:         u.Version,
			Source:          u.Source,
			CreatedAt:       u.CreatedAt,
			UpdatedAt:       u.UpdatedAt,
		}
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

	unit := models.Unit{
		ID:              dbUnit.ID,
		FactionID:       dbUnit.FactionID,
		Name:            dbUnit.Name,
		Description:     dbUnit.Description,
		IsManifestation: dbUnit.IsManifestation,
		Move:            int(dbUnit.Move),
		Health:          int(dbUnit.Health),
		Save:            dbUnit.Save,
		Ward:            dbUnit.Ward,
		Control:         int(dbUnit.Control),
		Points:          int(dbUnit.Points),
		SummonCost:      dbUnit.SummonCost,
		Banishment:      dbUnit.Banishment,
		MinSize:         int(dbUnit.MinSize),
		MaxSize:         int(dbUnit.MaxSize),
		MatchedPlay:     dbUnit.MatchedPlay,
		Version:         dbUnit.Version,
		Source:          dbUnit.Source,
		CreatedAt:       dbUnit.CreatedAt,
		UpdatedAt:       dbUnit.UpdatedAt,
	}

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
	if dbManifestations == nil {
		return []models.Unit{}, nil
	}

	manifestations := make([]models.Unit, len(dbManifestations))
	for i, m := range dbManifestations {
		manifestations[i] = models.Unit{
			ID:              m.ID,
			FactionID:       m.FactionID,
			Name:            m.Name,
			Description:     m.Description,
			IsManifestation: m.IsManifestation,
			Move:            int(m.Move),
			Health:          int(m.Health),
			Save:            m.Save,
			Ward:            m.Ward,
			Control:         int(m.Control),
			Points:          int(m.Points),
			SummonCost:      m.SummonCost,
			Banishment:      m.Banishment,
			MinSize:         int(m.MinSize),
			MaxSize:         int(m.MaxSize),
			MatchedPlay:     m.MatchedPlay,
			Version:         m.Version,
			Source:          m.Source,
			CreatedAt:       m.CreatedAt,
			UpdatedAt:       m.UpdatedAt,
		}
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
	if dbUnits == nil {
		return []models.Unit{}, nil
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = models.Unit{
			ID:              u.ID,
			FactionID:       u.FactionID,
			Name:            u.Name,
			Description:     u.Description,
			IsManifestation: u.IsManifestation,
			Move:            int(u.Move),
			Health:          int(u.Health),
			Save:            u.Save,
			Ward:            u.Ward,
			Control:         int(u.Control),
			Points:          int(u.Points),
			SummonCost:      u.SummonCost,
			Banishment:      u.Banishment,
			MinSize:         int(u.MinSize),
			MaxSize:         int(u.MaxSize),
			MatchedPlay:     u.MatchedPlay,
			Version:         u.Version,
			Source:          u.Source,
			CreatedAt:       u.CreatedAt,
			UpdatedAt:       u.UpdatedAt,
		}
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

	unit := models.Unit{
		ID:              dbManifestation.ID,
		FactionID:       dbManifestation.FactionID,
		Name:            dbManifestation.Name,
		Description:     dbManifestation.Description,
		IsManifestation: dbManifestation.IsManifestation,
		Move:            int(dbManifestation.Move),
		Health:          int(dbManifestation.Health),
		Save:            dbManifestation.Save,
		Ward:            dbManifestation.Ward,
		Control:         int(dbManifestation.Control),
		Points:          int(dbManifestation.Points),
		SummonCost:      dbManifestation.SummonCost,
		Banishment:      dbManifestation.Banishment,
		MinSize:         int(dbManifestation.MinSize),
		MaxSize:         int(dbManifestation.MaxSize),
		Version:         dbManifestation.Version,
		Source:          dbManifestation.Source,
		CreatedAt:       dbManifestation.CreatedAt,
		UpdatedAt:       dbManifestation.UpdatedAt,
	}

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

	if dbUnits == nil {
		return []models.Unit{}, nil
	}

	units := make([]models.Unit, len(dbUnits))
	for i, u := range dbUnits {
		units[i] = models.Unit{
			ID:              u.ID,
			FactionID:       u.FactionID,
			Name:            u.Name,
			Description:     u.Description,
			IsManifestation: u.IsManifestation,
			Move:            int(u.Move),
			Health:          int(u.Health),
			Save:            u.Save,
			Ward:            u.Ward,
			Control:         int(u.Control),
			Points:          int(u.Points),
			SummonCost:      u.SummonCost,
			Banishment:      u.Banishment,
			MinSize:         int(u.MinSize),
			MaxSize:         int(u.MaxSize),
			MatchedPlay:     u.MatchedPlay,
			Version:         u.Version,
			Source:          u.Source,
			CreatedAt:       u.CreatedAt,
			UpdatedAt:       u.UpdatedAt,
		}
	}

	return units, nil
}
