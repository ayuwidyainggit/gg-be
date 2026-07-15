package validation

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

)

func NewProductValidator() *Validate {
	v := NewValidator() // pakai base dari global validator

	// Daftarkan validasi baru
	v.Validator.RegisterValidation("alphanumericSpaceDash", alphanumericSpaceDash)

	// Tambahkan terjemahan untuk bahasa Inggris
	v.Validator.RegisterTranslation("alphanumericSpaceDash", transEn, func(ut ut.Translator) error {
		return ut.Add("alphanumericSpaceDashEN", "{0} must contain only letters, numbers, spaces, or dashes (-)", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericSpaceDashEN", fe.Field())
		return t
	})

	// Tambahkan terjemahan untuk bahasa Indonesia
	v.Validator.RegisterTranslation("alphanumericSpaceDash", transId, func(ut ut.Translator) error {
		return ut.Add("alphanumericSpaceDashID", "{0} hanya boleh berisi huruf, angka, spasi, atau tanda minus (-)", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericSpaceDashID", fe.Field())
		return t
	})

	return v
}

func alphanumericSpaceDash(fl validator.FieldLevel) bool {
	string := fl.Field().String()

	return regexp.MustCompile(`^[a-zA-Z0-9\s-]*$`).MatchString(string)
}
