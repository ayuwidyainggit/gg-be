package exception

type ForbiddenError struct {
	Error string
}

func NewForbiddenError(error string) ForbiddenError {
	return ForbiddenError{Error: error}
}
