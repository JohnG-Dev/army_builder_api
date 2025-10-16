package models

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID              uuid.UUID `json:"id"`
	FactionID       uuid.UUID `json:"faction_id"`
	Name            string    `json:"name"`
	IsManifestation bool      `json:"is_manifestation"`
	Move            int       `json:"move"`
	Wounds          int       `json:"wounds"`
	Save            string    `json:"save"`
	Ward            string    `json:"ward"`
	Control         int       `json:"control"`
	MinSize         int       `json:"min_size"`
	MaxSize         int       `json:"max_size"`
	Version         string    `json:"version"`
	Source          string    `json:"source"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
