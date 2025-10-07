package errors

import "errors"

var (
	ErrMissingUnitID = errors.New("Unit_ID paramater required")
)
