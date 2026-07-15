package validation

import (
	"fmt"
	"master/pkg/constant"
	"master/pkg/str"
	"reflect"
	"regexp"
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

func NewValidator() *Validate {
	vc := validator.New()
	uni := ut.New(BaseLocale, SupportedLocales.Translators()...)
	universalTranslators := make(LocaleUniversalTranslators, len(SupportedLocales))

	vc.RegisterValidation("qtystr", qtystr)
	vc.RegisterValidation("alphanumericSpace", alphanumericSpace)
	vc.RegisterValidation("alphanumDashUnderscore", alphanumDashUnderscore)
	vc.RegisterValidation("answer_frequency", answerFrequency)
	vc.RegisterValidation("level_target", levelTarget)

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
			vc.RegisterTranslation("qtystr", transId, func(ut ut.Translator) error {
				err := ut.Add("qtyStrID", "{0} nilai harus valid! contoh : 00001.000.000", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("qtyStrID", fe.Field())
				return t
			})

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
			vc.RegisterTranslation("alphanumDashUnderscore", transId, func(ut ut.Translator) error {
				err := ut.Add("alphanumDashUnderscoreID", "{0} harus alfanumerik, tanda hubung (-), atau garis bawah (_), maksimal 20 karakter", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumDashUnderscoreID", fe.Field())
				return t
			})
			vc.RegisterTranslation("answer_frequency", transId, func(ut ut.Translator) error {
				err := ut.Add("answerFrequencyID", "{0} harus salah satu dari: One Time, Multiple Times, One Day, Multiple Times, Different Day", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("answerFrequencyID", fe.Field())
				return t
			})
			vc.RegisterTranslation("level_target", transId, func(ut ut.Translator) error {
				err := ut.Add("levelTargetID", "{0} harus salah satu dari: Salesman, Outlet, Distributor", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("levelTargetID", fe.Field())
				return t
			})
		case en.New().Locale(): // an inefficiency bc locale names are not available as constants
			if err := enTranslations.RegisterDefaultTranslations(vc, utr); err != nil {
				panic(err)
			}
			transEn = utr
			vc.RegisterTranslation("qtystr", transEn, func(ut ut.Translator) error {
				err := ut.Add("qtyStrEN", "{0} must be a valid value! contoh : 00001.000.000", true) // see universal-translator for details
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("qtyStrEN", fe.Field())
				return t
			})

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
			vc.RegisterTranslation("alphanumDashUnderscore", transEn, func(ut ut.Translator) error {
				err := ut.Add("alphanumDashUnderscoreEN", "{0} must be alphanumeric, dash (-), or underscore (_), maximum 20 characters", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("alphanumDashUnderscoreEN", fe.Field())
				return t
			})
			vc.RegisterTranslation("answer_frequency", transEn, func(ut ut.Translator) error {
				err := ut.Add("answerFrequencyEN", "{0} must be one of: One Time, Multiple Times, One Day, Multiple Times, Different Day", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("answerFrequencyEN", fe.Field())
				return t
			})
			vc.RegisterTranslation("level_target", transEn, func(ut ut.Translator) error {
				err := ut.Add("levelTargetEN", "{0} must be one of: Salesman, Outlet, Distributor", true)
				if err != nil {
					return err
				}
				return nil
			}, func(ut ut.Translator, fe validator.FieldError) string {
				var t string
				t, _ = ut.T("levelTargetEN", fe.Field())
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

	v.Validator.RegisterTranslation("qtystr", v.Translator, func(ut ut.Translator) error {
		err := ut.Add("qtyStrEn", "{0} must be alphabet, number or space only", true) // see universal-translator for details
		if err != nil {
			return err
		}
		err = ut.Add("qtyStrID", "{0} harus alfabet, angka atau spasi saja", true) // see universal-translator for details
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

}

func (v *Validate) ValidateStruct(request interface{}, lang string) []map[string]interface{} {
	if lang == "id" {
		v.Translator = transId
	} else {
		v.Translator = transEn
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
			if len([]rune(parse[i])) != 3 {
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

// alphanumDashUnderscore accepts A-Z, a-z, 0-9, dash, and underscore; max 20 chars after trim.
// ponytail: superset dari bawaan `alphanum` (^[0-9a-zA-Z]*$) — tambahkan `-` dan `_`. Trimming
// dilakukan agar payload "  DIST001  " tidak lolos jadi beda key.
var rxDistributorCode = regexp.MustCompile(`^[A-Za-z0-9_-]{1,20}$`)

func alphanumDashUnderscore(fl validator.FieldLevel) bool {
	s := strings.TrimSpace(fl.Field().String())
	if s == "" {
		return false
	}
	return rxDistributorCode.MatchString(s)
}

func answerFrequency(fl validator.FieldLevel) bool {
	return constant.IsValidSurveyAnswerFrequencyForWrite(fl.Field().String())
}

// levelTarget accepts only the Sprint 13 enum values from the DOCX
// (Enhance_Create_Survey_BE, body field table). Empty input is intentionally
// accepted because the validator is paired with `omitempty` on entity
// CreateSurveyBody/UpdateSurveyBody to remain backward compatible with legacy
// payloads. Non-empty values that are not in the allow list must fail.
func levelTarget(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	switch value {
	case "Outlet", "Distributor", "Salesman":
		return true
	}
	return false
}
