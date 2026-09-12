package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

const MessageValidationFailed = "validation error"

// RegisterJSONTagNames makes validation error messages report JSON field names
// (e.g. "disaster_id") instead of Go struct field names (e.g. "DisasterID").
func RegisterJSONTagNames(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Message converts a binding/validation error into a human-readable message
// for the first failing field.
func Message(err error) string {
	details := Details(err)
	if len(details) == 0 {
		return "request body is invalid"
	}
	return details[0].Message
}

// Details converts a binding/validation error into one structured entry per
// failing field. Fields use their JSON tag path (e.g. "location.source").
func Details(err error) []apperror.Detail {
	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) {
		if isNumeric(typeError.Type) {
			return []apperror.Detail{{Field: typeError.Field, Message: fmt.Sprintf("%s must be a number", typeError.Field)}}
		}
		return []apperror.Detail{{Field: typeError.Field, Message: fmt.Sprintf("%s has an invalid value", typeError.Field)}}
	}

	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) || len(fieldErrors) == 0 {
		return nil
	}

	details := make([]apperror.Detail, 0, len(fieldErrors))
	for _, fieldError := range fieldErrors {
		details = append(details, apperror.Detail{
			Field:   jsonFieldPath(fieldError),
			Message: messageFor(fieldError),
		})
	}
	return details
}

func jsonFieldPath(fieldError validator.FieldError) string {
	ns := fieldError.Namespace()
	if i := strings.Index(ns, "."); i >= 0 {
		return ns[i+1:]
	}
	return ns
}

func messageFor(fieldError validator.FieldError) string {
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
