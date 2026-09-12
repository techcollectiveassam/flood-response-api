package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

type Params struct {
	DefaultPage  int
	DefaultLimit int
	MaxLimit     int
}

type Query struct {
	Page  *int `form:"page" binding:"omitempty,min=1"`
	Limit *int `form:"limit" binding:"omitempty,min=1"`
}

func ParamsFromConfig(cfg *config.Config) Params {
	if cfg == nil {
		return Params{}
	}
	return Params{
		DefaultPage:  cfg.Pagination.DefaultPage,
		DefaultLimit: cfg.Pagination.DefaultLimit,
		MaxLimit:     cfg.Pagination.MaxLimit,
	}
}

type contextKey struct{}

func Middleware(params Params) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(contextKey{}, params)
		c.Next()
	}
}

func FromContext(c *gin.Context) Params {
	params, ok := c.Get(contextKey{})
	if !ok {
		return Params{}
	}
	if res, ok := params.(Params); ok {
		return res
	}
	return Params{}
}

func Resolve(query Query, params Params) (int, int, error) {
	page, limit, max := params.DefaultPage, params.DefaultLimit, params.MaxLimit
	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if max < 1 {
		max = maxLimit
	}
	if max < limit {
		max = limit
	}

	if query.Page != nil {
		page = *query.Page
	}
	if query.Limit != nil {
		if *query.Limit > max {
			return 0, 0, apperror.BadRequest("invalid_query_params", "limit must not exceed "+strconv.Itoa(max))
		}
		limit = *query.Limit
	}

	return page, limit, nil
}

func TotalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}
	return int((total + int64(limit) - 1) / int64(limit))
}
