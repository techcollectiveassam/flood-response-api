package validation

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

type disasterRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required,oneof=flood earthquake landslide"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
}

type locationPayload struct {
	Source string `json:"source" binding:"required,oneof=gps place_search map text"`
}

type affectedAreaRequest struct {
	DisasterID int32            `json:"disaster_id" binding:"required"`
	Name       string           `json:"name" binding:"required"`
	Location   *locationPayload `json:"location" binding:"required"`
}

func newTestValidator(t *testing.T) *validator.Validate {
	t.Helper()
	validate := validator.New()
	validate.SetTagName("binding")
	RegisterJSONTagNames(validate)
	return validate
}

func TestMessageUsesJSONFieldNames(t *testing.T) {
	validate := newTestValidator(t)

	t.Run("required uses json field name", func(t *testing.T) {
		req := affectedAreaRequest{Name: "area", Location: &locationPayload{Source: "gps"}}
		err := validate.Struct(req)
		require.NotNil(t, err)
		assert.Equal(t, "disaster_id is required", Message(err))
	})

	t.Run("oneof uses json field name and params", func(t *testing.T) {
		req := disasterRequest{Name: "Flood", Type: "cyclone"}
		err := validate.Struct(req)
		require.NotNil(t, err)
		assert.Equal(t, "type must be one of: flood, earthquake, landslide", Message(err))
	})

	t.Run("missing nested struct reports its own name", func(t *testing.T) {
		req := affectedAreaRequest{DisasterID: 1, Name: "area"}
		err := validate.Struct(req)
		require.NotNil(t, err)
		assert.Equal(t, "location is required", Message(err))
	})

	t.Run("default branch uses json field name", func(t *testing.T) {
		req := locationPayload{Source: "invalid"}
		err := validate.Struct(req)
		require.NotNil(t, err)
		assert.Equal(t, "source must be one of: gps, place_search, map, text", Message(err))
	})
}

func TestDetailsReturnsAllErrorsWithJSONPaths(t *testing.T) {
	validate := newTestValidator(t)

	t.Run("all errors surfaced in order", func(t *testing.T) {
		req := affectedAreaRequest{Location: &locationPayload{Source: ""}}
		err := validate.Struct(req)
		require.NotNil(t, err)
		assert.Equal(t, []apperror.Detail{
			{Field: "disaster_id", Message: "disaster_id is required"},
			{Field: "name", Message: "name is required"},
			{Field: "location.source", Message: "source is required"},
		}, Details(err))
	})

	t.Run("missing nested struct uses its own path", func(t *testing.T) {
		req := affectedAreaRequest{DisasterID: 1, Name: "area"}
		err := validate.Struct(req)
		require.NotNil(t, err)
		assert.Equal(t, []apperror.Detail{
			{Field: "location", Message: "location is required"},
		}, Details(err))
	})
}

func TestMessageUnmarshalTypeError(t *testing.T) {
	err := &json.UnmarshalTypeError{Field: "disaster_id", Type: reflect.TypeOf(int32(0))}
	assert.Equal(t, "disaster_id must be a number", Message(err))
	assert.Equal(t, []apperror.Detail{{Field: "disaster_id", Message: "disaster_id must be a number"}}, Details(err))
}

func TestMessageWithoutTagNameFuncStillReadable(t *testing.T) {
	validate := validator.New()
	validate.SetTagName("binding")

	req := affectedAreaRequest{Name: "area", Location: &locationPayload{Source: "gps"}}
	err := validate.Struct(req)
	require.NotNil(t, err)
	assert.Equal(t, "DisasterID is required", Message(err))
}
