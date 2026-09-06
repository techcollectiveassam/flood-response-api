package affectedarea

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/affected-areas", m.handler.CreateAffectedArea)
}
