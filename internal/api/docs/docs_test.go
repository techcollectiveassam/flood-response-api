package docs_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/api/docs"
)

func TestRegisterServesRedocPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	docs.Register(router)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, w.Body.String(), "Flood Response API")
}

func TestRegisterServesRedocBundle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	docs.Register(router)

	req := httptest.NewRequest(http.MethodGet, "/docs/redoc.standalone.js", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/javascript")
	assert.NotEmpty(t, w.Body.Bytes())
}
