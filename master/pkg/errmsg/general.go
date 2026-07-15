package errmsg

const (
	ERROR_PARAM_REQUIRED       = "%s is required"
	ERROR_PARAM_MIN            = "%s minimum is %s"
	ERROR_PARAM_MAX            = "%s maximum is %s"
	ERROR_PARAM_INVALID        = "%s is invalid"
	ERROR_DATE_FORMAT          = "%s: %s, invalid date format"
	ERROR_DEL_DATA_NOT_ALLOWED = "This warehouse is not allowed to be deleted"
	ERROR_DATE_MUST_GT_NOW     = `must be more than current date + 1`
	ERROR_NO_ROWS_AFFECTED     = "no rows affected"
)
