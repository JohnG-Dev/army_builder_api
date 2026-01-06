package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/JohnG-Dev/army_builder_api/internal/database"
	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

type Seeder struct {
	s          *state.State
	ctx        context.Context
	keywordMap map[string]uuid.UUID
	gameMap    map[string]uuid.UUID
	txQueries  *database.Queries
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
	defer tx.Rollback(sr.ctx)

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

	gameName, err := sr.getDB().GetGameByName(sr.ctx, name)
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
	sr.gameMap[name] = gameName.ID
	return gameName.ID, nil
}

func (sr *Seeder) createFaction(gameID uuid.UUID, f models.FactionSeed) (uuid.UUID, error) {
	faction, err := sr.getDB().CreateFaction(sr.ctx, database.CreateFactionParams{
		GameID:     gameID,
		Name:       f.Name,
		Allegiance: f.Allegiance,
		Version:    f.Version,
		Source:     f.Source,
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
		Move:              cleanStat(u.Move),
		HealthWounds:      cleanStat(u.Health),
		SaveStats:         cleanStat(u.Save),
		WardFnp:           cleanStat(u.Ward),
		InvulnSave:        cleanStat(u.Invuln),
		ControlOc:         cleanStat(u.Control),
		Toughness:         cleanStat(u.Toughness),
		LeadershipBravery: cleanStat(u.Leadership),
		Points:            int32(u.Points),
		AdditionalStats:   statsJSON,
		SummonCost:        u.SummonCost,
		Banishment:        u.Banishment,
		MinUnitSize:       int32(u.MinUnitSize),
		MaxUnitSize:       int32(u.MaxUnitSize),
		MatchedPlay:       u.MatchedPlay,
		Version:           version,
		Source:            source,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return unit.ID, nil
}

func (sr *Seeder) seedUnitWeapons(unitID uuid.UUID, weapons []models.WeaponSeed, version, source string) error {
	sr.s.Logger.Info("Seeding Weapons", zap.Int("count", len(weapons)))
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
		sr.s.Logger.Info("Seeding Ability", zap.String("name", a.Name))
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
		if len(a.Effects) > 0 {
			sr.s.Logger.Info("Seeding Ability Effects", zap.Int("count", len(a.Effects)))
		}
		for _, e := range a.Effects {
			_, err := sr.getDB().CreateAbilityEffect(sr.ctx, database.CreateAbilityEffectParams{
				AbilityID:   ability.ID,
				Stat:        e.Stat,
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
		sr.s.Logger.Info("Seeding BattleFormations", zap.String("name", b.Name))
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
		sr.s.Logger.Info("Seeding Enhancement", zap.String("name", e.Name))
		_, err := sr.getDB().CreateEnhancement(sr.ctx, database.CreateEnhancementParams{
			FactionID:       factionID,
			Name:            e.Name,
			EnhancementType: e.EnhancementType,
			Description:     e.Description,
			Restrictions:    e.Restrictions,
			Points:          int32(e.Points),
			Version:         version,
			Source:          source,
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
