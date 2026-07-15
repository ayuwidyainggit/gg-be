package validation

import (
	"fmt"
	"inventory/pkg/str"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

var (
	// uni *ut.UniversalTranslator
	BaseLocale       locales.Translator = en.New()
	SupportedLocales LocaleTranslators  = LocaleTranslators{
		BaseLocale.Locale(): BaseLocale,
		id.New().Locale():   id.New(),
		// ^ an inefficiency bc locale names are not available as constants
	}
	transEn ut.Translator
	transId ut.Translator
)

type LocaleTranslators map[string]locales.Translator

func (lt LocaleTranslators) Translators() []locales.Translator {
	ts := make([]locales.Translator, 0, len(lt))
	for _, t := range lt {
		ts = append(ts, t)
	}
	return ts
}

type LocaleUniversalTranslators map[string]ut.Translator

type Validate struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() *Validate {
	vc := validator.New()
	uni := ut.New(BaseLocale, SupportedLocales.Translators()...)
	universalTranslators := make(LocaleUniversalTranslators, len(SupportedLocales))

	vc.RegisterValidation("yyyyMmDdDate", yyyyMmDdDate)
	vc.RegisterValidation("alphanum", alphanum)
	vc.RegisterValidation("alphanumspace", alphanumspace)
	vc.RegisterValidation("ddMmYyyyDate", ddMmYyyyDate)

	for l := range SupportedLocales {
		utr, ok := uni.GetTranslator(l)
		if !ok {
			panic(fmt.Sprintf("%s translator not found", l))
		}
		// an inefficiency bc translations/* are independent packages instead of a single package
		switch l {
		case id.New().Locale(): // an inefficiency bc locale names are not available as constants
			if err := idTranslations.RegisterDefaultTranslations(vc, utr); err != nil {
				panic(err)
			}
			transId = utr

			vc.RegisterTranslation("yyyyMmDdDate", transId, func(ut ut.Translator) error {
				err := ut.Add("yyyyMmDdDateID", "{0} harus yyyy-mm-dd", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("yyyyMmDdDateID", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanum", transId, func(ut ut.Translator) error {
				err := ut.Add("alphanumID", "{0} harus alphanumeric", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumID", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanumspace", transId, func(ut ut.Translator) error {
				err := ut.Add("alphanumspaceID", "{0} hanya boleh huruf, angka, dan spasi", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumspaceID", fe.Field())
				return t
			})

			vc.RegisterTranslation("ddMmYyyyDate", transId, func(ut ut.Translator) error {
				err := ut.Add("ddMmYyyyDateID", "{0} harus dd/mm/yyyy", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("ddMmYyyyDateID", fe.Field())
				return t
			})
		case en.New().Locale(): // an inefficiency bc locale names are not available as constants
			if err := enTranslations.RegisterDefaultTranslations(vc, utr); err != nil {
				panic(err)
			}
			transEn = utr

			vc.RegisterTranslation("yyyyMmDdDate", transEn, func(ut ut.Translator) error {
				err := ut.Add("yyyyMmDdDateEN", "{0} must be yyyy-mm-dd", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("yyyyMmDdDateEN", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanum", transEn, func(ut ut.Translator) error {
				err := ut.Add("alphanumEN", "{0} must be alphanumeric", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumEN", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanumspace", transEn, func(ut ut.Translator) error {
				err := ut.Add("alphanumspaceEN", "{0} must contain only letters, numbers, and spaces", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumspaceEN", fe.Field())
				return t
			})

			vc.RegisterTranslation("ddMmYyyyDate", transEn, func(ut ut.Translator) error {
				err := ut.Add("ddMmYyyyDateEN", "{0} must be dd/mm/yyyy", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("ddMmYyyyDateEN", fe.Field())
				return t
			})

		}

		universalTranslators[l] = utr
	}

	vc.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validate{
		Validator:  vc,
		Translator: transEn,
	}
}

func (v *Validate) ValidateStruct(request interface{}, lang string) []map[string]interface{} {
	if lang == "id" {
		v.Translator = transId
	}
	err := v.Validator.Struct(request)
	errs := v.translateError(err)
	return errs
}

func (v *Validate) ValidateStructReturnError(request interface{}) error {
	err := v.Validator.Struct(request)
	if err != nil {
		validatorErrs := err.(validator.ValidationErrors)
		errorTranslation := validatorErrs.Translate(v.Translator)
		for _, e := range validatorErrs {
			errorMsg := errorTranslation[e.Namespace()]
			err = fmt.Errorf(errorMsg)
			return err
		}
	}
	return nil
}

func (v *Validate) translateError(err error) (errors []map[string]interface{}) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	errors = make([]map[string]interface{}, 0)
	errorTranslation := validatorErrs.Translate(v.Translator)
	for _, e := range validatorErrs {
		errKey := e.Field()
		errorMsg := str.Replacer(errorTranslation[e.Namespace()], strings.NewReplacer(e.StructField(), errKey))
		newReplacerErrMsg := msgForTag(e.Tag())
		if newReplacerErrMsg != "" {
			errorMsg = newReplacerErrMsg
		}
		errorRow := map[string]interface{}{
			"key":     errKey,
			"message": errorMsg,
		}
		errors = append(errors, errorRow)
	}
	return errors
}

func msgForTag(tag string) string {
	switch tag {
	case "unique":
		return "This field must be unique"
	}
	return ""
}

func yyyyMmDdDate(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		return false
	}
	return true
}

func alphanum(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	// Allow empty string (use omitempty or required separately)
	if str == "" {
		return true
	}
	// Match alphanumeric characters only (letters and numbers)
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", str)
	return matched
}

func alphanumspace(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "" {
		return true
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9 ]+$", str)
	return matched
}

func ddMmYyyyDate(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	// Allow empty string (use omitempty or required separately)
	if str == "" {
		return true
	}
	// Parse DD/MM/YYYY format
	_, err := time.Parse("02/01/2006", str)
	if err != nil {
		return false
	}
	return true
}
