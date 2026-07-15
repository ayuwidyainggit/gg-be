package constant

import "errors"

var (
	ErrRecordNotFound              = errors.New("record not found")
	ErrStockOpnameCannotBeStarted  = errors.New("stock opname cannot be started")
	ErrNoRowsAffected              = errors.New("no rows affected")
	ErrDataNotFound              = errors.New("data not found")

)