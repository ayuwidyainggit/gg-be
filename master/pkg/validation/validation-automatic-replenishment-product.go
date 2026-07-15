package validation

func NewAutomaticReplenishmentProductValidator() *Validate {
	v := NewValidator() // use base from global validator

	// Register custom validations if needed
	// v.Validator.RegisterValidation("customValidation", customValidationFunc)

	// Add translations for custom validations if any
	// v.Validator.RegisterTranslation("customValidation", transEn, func(ut ut.Translator) error {
	// 	return ut.Add("customValidationEN", "{0} must be valid", true)
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("customValidationEN", fe.Field())
	// 	return t
	// })

	// v.Validator.RegisterTranslation("customValidation", transId, func(ut ut.Translator) error {
	// 	return ut.Add("customValidationID", "{0} harus valid", true)
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("customValidationID", fe.Field())
	// 	return t
	// })

	return v
}
