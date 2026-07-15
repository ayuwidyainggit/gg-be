package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"system/pkg/str"

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

	for l := range SupportedLocales {
		ut, ok := uni.GetTranslator(l)
		if !ok {
			panic(fmt.Sprintf("%s translator not found", l))
		}
		// an inefficiency bc translations/* are independent packages instead of a single package
		switch l {
		case id.New().Locale(): // an inefficiency bc locale names are not available as constants
			if err := idTranslations.RegisterDefaultTranslations(vc, ut); err != nil {
				panic(err)
			}
			transId = ut
		case en.New().Locale(): // an inefficiency bc locale names are not available as constants
			if err := enTranslations.RegisterDefaultTranslations(vc, ut); err != nil {
				panic(err)
			}
			transEn = ut
		}

		universalTranslators[l] = ut
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
