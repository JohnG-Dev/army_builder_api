package models

import (
	"time"

	"github.com/google/uuid"
)

// Rule represents a game-level core rule, battle tactic, or grand strategy.
type Rule struct {
	ID        uuid.UUID `json:"id"`
	GameID    uuid.UUID `json:"game_id"`
	Name      string    `json:"name"`
	RuleType  string    `json:"rule_type"` // core / battle_tactic / grand_strategy
	Text      string    `json:"text"`
	Version   string    `json:"version"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
