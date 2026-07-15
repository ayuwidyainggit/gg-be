package utils

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

func InitializeValidator() *validator.Validate {
	validate := validator.New()

	_ = validate.RegisterValidation("notEmptyStringSlice", func(fl validator.FieldLevel) bool {
		slices := fl.Field().Interface().([]string)
		if len(slices) == 0 {
			return false
		}
		for _, s := range slices {
			if s == "" {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		dateRegex := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
		if !dateRegex.MatchString(fl.Field().String()) {
			return false
		}
		return true
	})

	_ = validate.RegisterValidation("notEmptyIntSlice", func(fl validator.FieldLevel) bool {
		slices := fl.Field().Interface().([]int)
		if len(slices) == 0 {
			return false
		}
		for _, val := range slices {
			if val == 0 {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("isString", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return true
	})

	_ = validate.RegisterValidation("isInt", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(int)
		if !ok {
			return false
		}

		return true
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return name
		}
		return name
	})

	return validate
}
