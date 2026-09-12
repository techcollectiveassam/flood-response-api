package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/techcollectiveassam/flood-response-api/internal/api/middleware"
	"github.com/techcollectiveassam/flood-response-api/internal/app"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/validation"
)

func NewRouter(features *app.Features, logger *slog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(middleware.LoggingMiddleware(logger), gin.Recovery())
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.RegisterJSONTagNames(v)
	}
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	features.AffectedArea.RegisterRoutes(v1)
	features.Disaster.RegisterRoutes(v1)

	return router
}
