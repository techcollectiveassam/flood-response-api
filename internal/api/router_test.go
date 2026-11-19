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
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

func newTestRouter(t *testing.T, allowedOrigins ...string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	logger := slog.Default()
	features := app.NewFeatures(&app.Application{
		Logger: logger,
		Config: &config.Config{},
	})
	return NewRouter(features, logger, allowedOrigins)
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

func TestCORSAllowsConfiguredOriginOnPreflight(t *testing.T) {
	router := newTestRouter(t, "https://app.example.org")

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/sos", nil)
	req.Header.Set("Origin", "https://app.example.org")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://app.example.org", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSHeadersOnSimpleRequest(t *testing.T) {
	router := newTestRouter(t, "https://app.example.org")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://app.example.org")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://app.example.org", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSBlocksUnlistedOrigin(t *testing.T) {
	router := newTestRouter(t, "https://app.example.org")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSBlocksAllOriginsWhenNoneConfigured(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://app.example.org")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestRouterRegistersAPIRoutes(t *testing.T) {
	router := newTestRouter(t)

	routes := make(map[string]bool)
	for _, r := range router.Routes() {
		routes[r.Method+" "+r.Path] = true
	}

	assert.True(t, routes["GET /health"])
	assert.True(t, routes["POST /api/v1/affected-areas"])
	assert.True(t, routes["POST /api/v1/command-centers"])
	assert.True(t, routes["POST /api/v1/disasters"])
	assert.True(t, routes["GET /api/v1/disasters"])
	assert.True(t, routes["GET /api/v1/disasters/:id"])
	assert.True(t, routes["GET /api/v1/disasters/:id/command-center"])
}
