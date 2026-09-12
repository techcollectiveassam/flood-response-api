package pagination

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

const (
	testDefaultPage  = 1
	testDefaultLimit = 20
	testMaxLimit     = 100
)

func testParams() Params {
	return Params{
		DefaultPage:  testDefaultPage,
		DefaultLimit: testDefaultLimit,
		MaxLimit:     testMaxLimit,
	}
}

func TestResolveDefaults(t *testing.T) {
	page, limit, err := Resolve(Query{}, testParams())

	assert.NoError(t, err)
	assert.Equal(t, testDefaultPage, page)
	assert.Equal(t, testDefaultLimit, limit)
}

func TestResolveExplicitValues(t *testing.T) {
	params := testParams()
	query := Query{
		Page:  intPointer(2),
		Limit: intPointer(10),
	}

	page, limit, err := Resolve(query, params)

	assert.NoError(t, err)
	assert.Equal(t, 2, page)
	assert.Equal(t, 10, limit)
}

func TestResolveRejectsLimitAboveMax(t *testing.T) {
	query := Query{Limit: intPointer(testMaxLimit + 1)}

	page, limit, err := Resolve(query, testParams())

	assert.Error(t, err)
	assert.Equal(t, 0, page)
	assert.Equal(t, 0, limit)
	assert.Equal(t, 400, apperror.HTTPStatus(err))
	assert.Equal(t, "invalid_query_params", apperror.Code(err))
}

func TestResolveFallsBackWhenParamsZero(t *testing.T) {
	page, limit, err := Resolve(Query{}, Params{})

	assert.NoError(t, err)
	assert.Equal(t, defaultPage, page)
	assert.Equal(t, defaultLimit, limit)
}

func TestParamsFromConfig(t *testing.T) {
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPage:  1,
			DefaultLimit: 50,
			MaxLimit:     200,
		},
	}

	params := ParamsFromConfig(cfg)

	assert.Equal(t, 1, params.DefaultPage)
	assert.Equal(t, 50, params.DefaultLimit)
	assert.Equal(t, 200, params.MaxLimit)
}

func TestParamsFromConfigNil(t *testing.T) {
	params := ParamsFromConfig(nil)

	assert.Equal(t, Params{}, params)
}

func TestMiddlewareInjectsParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", Middleware(testParams()), func(c *gin.Context) {
		params := FromContext(c)
		c.JSON(200, params)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got Params
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, testParams(), got)
}

func TestFromContextWithoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		params := FromContext(c)
		c.JSON(200, params)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got Params
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, Params{}, got)
}

func TestTotalPages(t *testing.T) {
	tests := []struct {
		total int64
		limit int
		want  int
	}{
		{total: 0, limit: 20, want: 0},
		{total: 20, limit: 20, want: 1},
		{total: 21, limit: 20, want: 2},
		{total: 57, limit: 20, want: 3},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, TotalPages(tt.total, tt.limit))
	}
}

func intPointer(value int) *int {
	return &value
}
