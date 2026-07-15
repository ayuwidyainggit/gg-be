package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"sales/pkg/str"
	"strconv"
	"strings"

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

func NewValiditor() *Validate {
	vc := validator.New()
	uni := ut.New(BaseLocale, SupportedLocales.Translators()...)
	universalTranslators := make(LocaleUniversalTranslators, len(SupportedLocales))

	vc.RegisterValidation("alphanumericSpace", alphanumericSpace)
	vc.RegisterValidation("alphanumericSlash", alphanumericSlash)
	vc.RegisterValidation("alphanumericSpaceSlashPercent", alphanumericSpaceSlashPercent)
	vc.RegisterValidation("budgetID", budgetID)

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

			vc.RegisterTranslation("alphanumericSpace", transId, func(ut ut.Translator) error {
				err := ut.Add("alphanumericSpaceID", "{0} harus alfabet, angka atau spasi saja", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumericSpaceID", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanumericSlash", transId, func(ut ut.Translator) error {
				err := ut.Add("alphanumericSlashID", "{0} harus alfabet, angka atau slash saja", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumericSlashID", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanumericSpaceSlashPercent", transId, func(ut ut.Translator) error {
				err := ut.Add("alphanumericSpaceSlashPercentID", "{0} harus alfabet, angka, spasi, slash atau persen saja", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumericSpaceSlashPercentID", fe.Field())
				return t
			})

			vc.RegisterTranslation("budgetID", transId, func(ut ut.Translator) error {
				return ut.Add("budgetIDID", "{0} only allowed letters, numbers, underscores, or hyphens", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("budgetIDID", fe.Field())
				return t
			})
		case en.New().Locale(): // an inefficiency bc locale names are not available as constants
			if err := enTranslations.RegisterDefaultTranslations(vc, utr); err != nil {
				panic(err)
			}
			transEn = utr

			vc.RegisterTranslation("alphanumericSpace", transEn, func(ut ut.Translator) error {
				err := ut.Add("alphanumericSpaceEN", "{0} must be alphabet, number or space only", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumericSpaceEN", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanumericSlash", transEn, func(ut ut.Translator) error {
				err := ut.Add("alphanumericSlashEN", "{0} must be alphabet, number or slash only", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumericSlashEN", fe.Field())
				return t
			})

			vc.RegisterTranslation("alphanumericSpaceSlashPercent", transEn, func(ut ut.Translator) error {
				err := ut.Add("alphanumericSpaceSlashPercentEN", "{0} must be alphabet, number, space, slash or percent only", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumericSpaceSlashPercentEN", fe.Field())
				return t
			})

			vc.RegisterTranslation("budgetID", transEn, func(ut ut.Translator) error {
				return ut.Add("budgetIDEN", "{0} must contain only letters, numbers, underscores, or hyphens", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("budgetIDEN", fe.Field())
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

func (v *Validate) RegisterCustomValidation() {
	v.Validator.RegisterValidation("qtystr", qtystr)
	v.Validator.RegisterTranslation("qtystr", v.Translator, func(ut ut.Translator) error {
		err := ut.Add("qtyStrEn", "{0} must be a valid value! example : 00001.000.000", true) // see universal-translator for details
		if err != nil {
			return err
		}
		err = ut.Add("qtyStrID", "{0} nilai harus valid! contoh : 00001.000.000", true) // see universal-translator for details
		if err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validator.FieldError) string {
		var t string
		switch ut.Locale() {
		case "en":
			t, _ = ut.T("qtyStrEn", fe.Field())
		case "id":
			t, _ = ut.T("qtyStrID", fe.Field())
		}

		return t
	})

	v.Validator.RegisterTranslation("required_with", v.Translator, func(ut ut.Translator) error {
		err := ut.Add("required_withEn", "{0} is required", true) // see universal-translator for details
		if err != nil {
			return err
		}
		err = ut.Add("required_withID", "{0} Wajib diisi", true) // see universal-translator for details
		if err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validator.FieldError) string {
		var t string
		switch ut.Locale() {
		case "en":
			t, _ = ut.T("required_withEn", fe.Field())
		case "id":
			t, _ = ut.T("required_withID", fe.Field())
		}

		return t
	})
}
func (v *Validate) ValidateStruct(request interface{}, lang string) []map[string]interface{} {
	if lang == "id" {
		v.Translator = transId
	}
	err := v.Validator.Struct(request)
	errs := v.translateError(err)
	return errs
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
		errorRow := map[string]interface{}{
			"key":     errKey,
			"message": errorMsg,
		}
		errors = append(errors, errorRow)
	}
	return errors
}

func qtystr(fl validator.FieldLevel) bool {
	qtyStr := fl.Field().String()
	parse := strings.Split(qtyStr, ".")
	if len(parse) != 3 {
		return false
	}
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			if len([]rune(parse[i])) != 5 {
				return false
			}

		case 1:
			if len([]rune(parse[i])) != 3 {
				return false
			}
		case 2:
			if len([]rune(parse[i])) != 2 {
				return false
			}
		}
		_, err := strconv.Atoi(parse[i])
		if err != nil {
			return false
		}
	}
	return true
}

func alphanumericSpace(fl validator.FieldLevel) bool {
	string := fl.Field().String()

	return regexp.MustCompile(`^[a-zA-Z0-9\s]*$`).MatchString(string)
}

func alphanumericSlash(fl validator.FieldLevel) bool {
	string := fl.Field().String()

	return regexp.MustCompile(`^[a-zA-Z0-9\/]*$`).MatchString(string)
}

func alphanumericSpaceSlashPercent(fl validator.FieldLevel) bool {
	string := fl.Field().String()

	return regexp.MustCompile(`^[a-zA-Z0-9\/\s\%]*$`).MatchString(string)
}

func budgetID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value)
}
