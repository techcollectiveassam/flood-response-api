package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/app"
)

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	features := app.NewFeatures(&app.Application{Logger: slog.Default()})
	return NewRouter(features)
}

func TestHealthEndpoint(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

func TestRouterRegistersAPIRoutes(t *testing.T) {
	router := newTestRouter(t)

	routes := make(map[string]bool)
	for _, r := range router.Routes() {
		routes[r.Method+" "+r.Path] = true
	}

	assert.True(t, routes["GET /health"])
	assert.True(t, routes["POST /api/v1/affected-areas"])
	assert.True(t, routes["POST /api/v1/disasters"])
	assert.True(t, routes["GET /api/v1/disasters"])
	assert.True(t, routes["GET /api/v1/disasters/:id"])
}
