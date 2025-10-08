package errors

import "errors"

var (
	//ErrMissingUnitID is returned when a weapon request is missing a unit_id
	ErrMissingUnitID = errors.New("unit_id parameter required")
)
