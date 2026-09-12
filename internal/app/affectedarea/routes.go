package affectedarea

import (
	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/pagination"
)

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/affected-areas", m.handler.CreateAffectedArea)
	rg.GET("/affected-areas", pagination.Middleware(m.params), m.handler.ListAffectedAreas)
}
