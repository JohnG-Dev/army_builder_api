package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type Seeder struct {
	s            *state.State
	ctx          context.Context
	keywordMap   map[string]uuid.UUID
	gameMap      map[string]uuid.UUID
	factionMap   map[string]uuid.UUID
	pendingLinks map[uuid.UUID]string
	txQueries    *database.Queries
}

func NewSeeder(ctx context.Context, s *state.State) *Seeder {
	return &Seeder{
		s:            s,
		ctx:          ctx,
		keywordMap:   make(map[string]uuid.UUID),
		gameMap:      make(map[string]uuid.UUID),
		factionMap:   make(map[string]uuid.UUID),
		pendingLinks: make(map[uuid.UUID]string),
	}
}

func (sr *Seeder) SeedFile(path string) error {
	// #nosec G304
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	var seed models.SeedData
	if err := yaml.Unmarshal(data, &seed); err != nil {
		return err
	}

	sr.s.Logger.Info("Processing file",
		zap.String("path", path),
		zap.String("game", seed.GameName),
	)

	tx, err := sr.s.Pool.Begin(sr.ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(sr.ctx) }()

	sr.txQueries = sr.s.DB.WithTx(tx)

	gameID, err := sr.getOrCreateGame(seed.GameName)
	if err != nil {
		return err
	}

	for _, f := range seed.Factions {
		sr.s.Logger.Info("Seeding Faction", zap.String("name", f.Name))
		factionID, err := sr.createFaction(gameID, f)
		if err != nil {
			return fmt.Errorf("failed to create faction %s: %w", f.Name, err)
		}

		sr.factionMap[f.Name] = factionID
		if f.ParentFactionName != "" {
			sr.pendingLinks[factionID] = f.ParentFactionName
		}

		err = sr.seedFactionBattleFormations(gameID, factionID, f.BattleFormations, f.Version, f.Source)
		if err != nil {
			return err
		}

		err = sr.seedFactionEnhancements(factionID, f.Enhancements, f.Version, f.Source)
		if err != nil {
			return err
		}

		for _, u := range f.Units {
			sr.s.Logger.Info("Seeding Unit", zap.String("name", u.Name))

			unitID, err := sr.createUnit(factionID, u, f.Version, f.Source)
			if err != nil {
				return fmt.Errorf("failed to create unit %s: %w", u.Name, err)
			}

			err = sr.seedUnitWeapons(unitID, u.Weapons, f.Version, f.Source)
			if err != nil {
				return err
			}

			err = sr.seedUnitKeywords(unitID, gameID, u.Keywords, f.Version, f.Source)
			if err != nil {
				return err
			}

			err = sr.seedUnitAbilities(unitID, factionID, gameID, u.Abilities, f.Version, f.Source)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit(sr.ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	sr.txQueries = nil
	return nil
}

func (sr *Seeder) LinkParents() error {
	sr.s.Logger.Info("Linking Faction Parents", zap.Int("count", len(sr.pendingLinks)))
	for factionID, parentName := range sr.pendingLinks {
		parentID, exists := sr.factionMap[parentName]
		if !exists {
			sr.s.Logger.Warn("Parent faction not found", zap.String("parent", parentName), zap.String("child_id", factionID.String()))
			continue
		}

		err := sr.s.DB.UpdateFactionParent(sr.ctx, database.UpdateFactionParentParams{
			ID:              factionID,
			ParentFactionID: database.UUIDToNullUUID(parentID),
			IsArmyOfRenown:  true,
		})
		if err != nil {
			return fmt.Errorf("failed to link parent for faction %s: %w", factionID, err)
		}
		sr.s.Logger.Info("Linked Parent", zap.String("child_id", factionID.String()), zap.String("parent_id", parentID.String()))
	}
	return nil
}

func (sr *Seeder) getDB() *database.Queries {
	if sr.txQueries != nil {
		return sr.txQueries
	}
	return sr.s.DB
}

func (sr *Seeder) getOrCreateGame(name string) (uuid.UUID, error) {
	id, exists := sr.gameMap[name]
	if exists {
		return id, nil
	}

	game, err := sr.getDB().GetGameByName(sr.ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			newGame, err := sr.s.DB.CreateGame(sr.ctx, database.CreateGameParams{
				Name:    name,
				Edition: "4th Edition",
				Version: "1.0",
				Source:  "Initial Seed",
			})
			if err != nil {
				return uuid.Nil, fmt.Errorf("failed to create game: %w", err)
			}
			sr.gameMap[name] = newGame.ID
			return newGame.ID, nil
		}
		return uuid.Nil, err
	}
	sr.gameMap[name] = game.ID
	return game.ID, nil
}

func (sr *Seeder) createFaction(gameID uuid.UUID, f models.FactionSeed) (uuid.UUID, error) {
	faction, err := sr.getDB().CreateFaction(sr.ctx, database.CreateFactionParams{
		GameID:             gameID,
		Name:               f.Name,
		Allegiance:         f.Allegiance,
		Version:            f.Version,
		Source:             f.Source,
		IsArmyOfRenown:     f.IsArmyOfRenown,
		IsRegimentOfRenown: f.IsRegimentOfRenown,
		ParentFactionID:    uuid.NullUUID{},
	})
	if err != nil {
		return uuid.Nil, err
	}

	return faction.ID, nil
}

func (sr *Seeder) createUnit(factionID uuid.UUID, u models.UnitSeed, version, source string) (uuid.UUID, error) {
	statsJSON, err := json.Marshal(u.AdditionalStats)
	if err != nil {
		statsJSON = []byte("{}")
	}
	unit, err := sr.getDB().CreateUnit(sr.ctx, database.CreateUnitParams{
		FactionID:         factionID,
		Name:              u.Name,
		IsUnique:          u.IsUnique,
		Move:              cleanStat(u.Move),
		HealthWounds:      cleanStat(u.Health),
		SaveStats:         cleanStat(u.Save),
		WardFnp:           cleanStat(u.Ward),
		InvulnSave:        cleanStat(u.Invuln),
		ControlOc:         cleanStat(u.Control),
		Toughness:         cleanStat(u.Toughness),
		LeadershipBravery: cleanStat(u.Leadership),
		// #nosec G115
		Points:          int32(u.Points),
		AdditionalStats: statsJSON,
		SummonCost:      u.SummonCost,
		Banishment:      u.Banishment,
		// #nosec G115
		MinUnitSize: int32(u.MinUnitSize),
		// #nosec G115
		MaxUnitSize: int32(u.MaxUnitSize),
		MatchedPlay: u.MatchedPlay,
		Version:     version,
		Source:      source,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return unit.ID, nil
}

func (sr *Seeder) seedUnitWeapons(unitID uuid.UUID, weapons []models.WeaponSeed, version, source string) error {
	for _, w := range weapons {
		_, err := sr.getDB().CreateWeapon(sr.ctx, database.CreateWeaponParams{
			UnitID:        unitID,
			Name:          w.Name,
			Range:         w.Range,
			Attacks:       w.Attacks,
			HitStats:      w.ToHit,
			WoundStrength: w.ToWound,
			RendAp:        w.Rend,
			Damage:        w.Damage,
			Version:       version,
			Source:        source,
		})
		if err != nil {
			return fmt.Errorf("failed to create weapon %s: %w", w.Name, err)
		}
	}

	return nil
}

func (sr *Seeder) seedUnitKeywords(unitID uuid.UUID, gameID uuid.UUID, keywordNames []string, version, source string) error {
	for _, name := range keywordNames {

		keywordID, exists := sr.keywordMap[name]

		if !exists {
			sr.s.Logger.Info("Creating New Keyword", zap.String("name", name))
			k, err := sr.getDB().CreateKeyword(sr.ctx, database.CreateKeywordParams{
				GameID:      gameID,
				Name:        name,
				Description: "",
				Version:     version,
				Source:      source,
			})
			if err != nil {
				return fmt.Errorf("failed to create keyword %s: %w", name, err)
			}
			keywordID = k.ID
			sr.keywordMap[name] = keywordID
		}
		err := sr.getDB().AddKeywordToUnit(sr.ctx, database.AddKeywordToUnitParams{
			UnitID:    unitID,
			KeywordID: keywordID,
			Value:     "",
		})
		if err != nil {
			return fmt.Errorf("failed to add keyword %s to unit: %w", name, err)
		}
	}
	return nil
}

func (sr *Seeder) seedUnitAbilities(unitID, factionID, gameID uuid.UUID, abilities []models.AbilitySeed, version, source string) error {
	for _, a := range abilities {
		ability, err := sr.getDB().CreateAbility(sr.ctx, database.CreateAbilityParams{
			UnitID:      database.UUIDToNullUUID(unitID),
			FactionID:   uuid.NullUUID{},
			GameID:      uuid.NullUUID{},
			Name:        a.Name,
			Description: a.Description,
			Type:        a.Type,
			Phase:       a.Phase,
			Version:     version,
			Source:      source,
		})
		if err != nil {
			return fmt.Errorf("failed to create ability %s: %w", a.Name, err)
		}
		for _, e := range a.Effects {
			_, err := sr.getDB().CreateAbilityEffect(sr.ctx, database.CreateAbilityEffectParams{
				AbilityID: ability.ID,
				Stat:      e.Stat,
				// #nosec G115
				Modifier:    int32(e.Modifier),
				Condition:   e.Condition,
				Description: e.Description,
				Version:     version,
				Source:      source,
			})
			if err != nil {
				return fmt.Errorf("failed to create ability effect %s: %w", a.Name, err)
			}
		}
	}

	return nil
}

func (sr *Seeder) seedFactionBattleFormations(gameID, factionID uuid.UUID, battleFormations []models.BattleFormationSeed, version, source string) error {
	for _, b := range battleFormations {
		_, err := sr.getDB().CreateBattleFormation(sr.ctx, database.CreateBattleFormationParams{
			GameID:      gameID,
			FactionID:   factionID,
			Name:        b.Name,
			Description: b.Description,
			Version:     version,
			Source:      source,
		})
		if err != nil {
			return fmt.Errorf("failed to create battleformation %s: %w", b.Name, err)
		}
	}
	return nil
}

func (sr *Seeder) seedFactionEnhancements(factionID uuid.UUID, enhancements []models.EnhancementSeed, version, source string) error {
	for _, e := range enhancements {
		_, err := sr.getDB().CreateEnhancement(sr.ctx, database.CreateEnhancementParams{
			FactionID:       factionID,
			Name:            e.Name,
			EnhancementType: e.EnhancementType,
			Description:     e.Description,
			Restrictions:    e.Restrictions,
			// #nosec G115
			Points:  int32(e.Points),
			Version: version,
			Source:  source,
		})
		if err != nil {
			return fmt.Errorf("failed to create enhancement %s: %w", e.Name, err)
		}
	}

	return nil
}

func cleanStat(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
