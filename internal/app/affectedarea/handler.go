package affectedarea

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAffectedArea(c *gin.Context) {
	var req CreateAffectedAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("invalid_request_body", validationMessage(err)))
		return
	}

	resp, err := h.service.CreateAffectedArea(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, resp)
}

func validationMessage(err error) string {
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
