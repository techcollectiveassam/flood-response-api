package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/app"
)

func NewRouter(features *app.Features) *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	features.AffectedArea.RegisterRoutes(v1)

	return router
}
