package errors

import "errors"

var (
	//ErrMissingUnitID is returned when a weapon request is missing a unit_id
	ErrMissingID        = errors.New("id paramater required")
	ErrMissingUnitID    = errors.New("unit id parameter required")
	ErrMissingFactionID = errors.New("faction id paramater required")
	ErrNotFound         = errors.New("resource not found")
)
