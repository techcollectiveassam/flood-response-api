package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Message converts a binding/validation error into a human-readable message.
func Message(err error) string {
	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) {
		if isNumeric(typeError.Type) {
			return fmt.Sprintf("%s must be a number", typeError.Field)
		}
		return fmt.Sprintf("%s has an invalid value", typeError.Field)
	}

	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) || len(fieldErrors) == 0 {
		return "request body is invalid"
	}

	fieldError := fieldErrors[0]
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldError.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fieldError.Field(), strings.Join(strings.Fields(fieldError.Param()), ", "))
	default:
		return fmt.Sprintf("%s is invalid", fieldError.Field())
	}
}

func isNumeric(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
