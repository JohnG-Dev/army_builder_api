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

func GetRulesForGame(s *state.State, ctx context.Context, gameID uuid.UUID) ([]models.Rule, error) {

	if gameID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	dbRules, err := s.DB.GetRulesForGame(ctx, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Rule{}, nil
		}

		return nil, err
	}

	if dbRules == nil {
		return []models.Rule{}, nil
	}

	rules := make([]models.Rule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = models.Rule{
			ID:        r.ID,
			GameID:    r.GameID,
			Name:      r.Name,
			RuleType:  r.RuleType,
			Text:      r.Text,
			Version:   r.Version,
			Source:    r.Source,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}

	return rules, nil
}

func GetRulesByType(s *state.State, ctx context.Context, gameID uuid.UUID, ruleType string) ([]models.Rule, error) {

	if gameID == uuid.Nil {
		return nil, appErr.ErrMissingID
	}

	if ruleType == "" {
		return nil, appErr.ErrMissingID
	}

	dbRules, err := s.DB.GetRulesByType(ctx, database.GetRulesByTypeParams{
		GameID:   gameID,
		RuleType: ruleType,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Rule{}, nil
		}

		return nil, err
	}

	if dbRules == nil {
		return []models.Rule{}, nil
	}

	rules := make([]models.Rule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = models.Rule{
			ID:        r.ID,
			GameID:    r.GameID,
			Name:      r.Name,
			RuleType:  r.RuleType,
			Text:      r.Text,
			Version:   r.Version,
			Source:    r.Source,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}

	return rules, nil
}

func GetRuleByID(s *state.State, ctx context.Context, id uuid.UUID) (models.Rule, error) {
	if id == uuid.Nil {
		return models.Rule{}, appErr.ErrMissingID
	}

	dbRule, err := s.DB.GetRuleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Rule{}, appErr.ErrNotFound
		}

		return models.Rule{}, err
	}

	rule := models.Rule{
		ID:        dbRule.ID,
		GameID:    dbRule.GameID,
		Name:      dbRule.Name,
		RuleType:  dbRule.RuleType,
		Text:      dbRule.Text,
		Version:   dbRule.Version,
		Source:    dbRule.Source,
		CreatedAt: dbRule.CreatedAt,
		UpdatedAt: dbRule.UpdatedAt,
	}

	return rule, nil
}
