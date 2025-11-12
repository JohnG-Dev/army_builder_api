package models

import (
	"time"

	"github.com/google/uuid"
)

// Keyword is the public tag vocabulary for units or manifestations.
type Keyword struct {
	ID          uuid.UUID `json:"id"`
	GameID      uuid.UUID `json:"game_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UnitKeyword links a keyword and an optional value to a unit (e.g., WARD 5+).
type UnitKeyword struct {
	UnitID    uuid.UUID `json:"unit_id"`
	KeywordID uuid.UUID `json:"keyword_id"`
	Value     string    `json:"value"`
}
