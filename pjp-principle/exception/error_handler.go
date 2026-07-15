package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"scyllax-pjp/data/response"
	"strings"
	"unicode"
)

func ErrorHandler(ctx *gin.Context, err interface{}) {
	if validationError(ctx, err) {
		return
	}

	if notFoundError(ctx, err) {
		return
	}

	if unauthorizedError(ctx, err) {
		return
	}

	if forbiddenError(ctx, err) {
		return
	}

	internalServerError(ctx, err)
}

func forbiddenError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(ForbiddenError)
	if ok {
		ctx.JSON(http.StatusForbidden, response.Error{
			Code:    http.StatusForbidden,
			Status:  "FORBIDDEN",
			Errors:  exception.Error,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func unauthorizedError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(UnauthorizedError)
	if ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Errors:  exception.Error,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func validationError(ctx *gin.Context, err interface{}) bool {

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		report := make(map[string]string)
		var fieldName string
		for _, e := range castedObject {
			if len(e.Namespace()) > 0 && unicode.IsUpper(rune(e.Namespace()[0])) {
				dotIndex := strings.Index(e.Namespace(), ".")
				if dotIndex != -1 {
					fieldName = e.Namespace()[dotIndex+1:]
				}
			} else {
				fieldName = e.Field()
			}
			switch e.Tag() {
			case "required":
				report[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				report[fieldName] = fmt.Sprintf("%s is not valid email", fieldName)
			case "gte":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "lte":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "unique":
				report[fieldName] = fmt.Sprintf("%s has already been taken", fieldName)
			case "max":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "min":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "numeric":
				report[fieldName] = fmt.Sprintf("%s value must be numeric", fieldName)
			case "digitNumeric":
				report[fieldName] = fmt.Sprintf("%s value must be numeric 4 digit", fieldName)
			case "oneof":
				report[fieldName] = fmt.Sprintf("%s value must be %s", fieldName, e.Param())
			case "len":
				report[fieldName] = fmt.Sprintf("%s value must be exactly %s characters long", fieldName, e.Param())
			case "notEmptyStringSlice":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is string", fieldName)
			case "dive":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty", fieldName)
			case "date":
				report[fieldName] = fmt.Sprintf("%s value must be date (yyyy-mm-dd)", fieldName)
			case "notEmptyIntSlice":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is int", fieldName)
			}
		}

		ctx.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  report,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func notFoundError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		ctx.JSON(http.StatusNotFound, response.Error{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Errors:  exception.Error,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func internalServerError(ctx *gin.Context, err interface{}) {
	var errMsg string
	if err != nil {
		errMsg = fmt.Sprintf("%v", err)
	} else {
		errMsg = "Unknown error occurred"
	}

	ctx.JSON(http.StatusInternalServerError, response.Error{
		Code:    http.StatusInternalServerError,
		Status:  "INTERNAL SERVER ERROR",
		Errors:  errMsg,
		TraceID: uuid.New().String(),
	})
}
