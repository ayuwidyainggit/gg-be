package apperr

type AppError struct {
	Code string
	Msg  string
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Msg
}

func New(code, msg string) *AppError {
	return &AppError{Code: code, Msg: msg}
}
