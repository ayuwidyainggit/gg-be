package validation

func NewProductAssignmentValidator() *Validate {
	v := NewValidator() // pakai base dari global validator

	// Add any custom validations if needed

	return v
}