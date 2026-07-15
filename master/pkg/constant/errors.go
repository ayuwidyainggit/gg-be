package constant

import "errors"

var (
	// ErrNoRowsAffected indicates that a database operation affected zero rows
	ErrNoRowsAffected = errors.New("no rows affected")

	// ErrSalesTargetNotFound indicates that the requested sales target does not exist
	ErrSalesTargetNotFound = errors.New("sales target not found")
)
