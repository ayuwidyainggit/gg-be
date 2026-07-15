package utils

import (
	"reflect"
	"regexp"
	"scyllax-pjp/helper"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeValidator(database *gorm.DB) *validator.Validate {
	db = database
	validate := validator.New()

	_ = validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		return helper.ValidateUnique(db, fl)
	})

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

	_ = validate.RegisterValidation("digitNumeric", func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface().(int)
		strVal := strconv.Itoa(val)
		if len(strVal) != 4 {
			return false
		}
		for _, r := range strVal {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		dateRegex := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
		if !dateRegex.MatchString(fl.Field().String()) {
			return false
		}
		return true
	})

	return validate
}
