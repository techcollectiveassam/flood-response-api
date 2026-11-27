package commandcenter

import (
	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/command-centers", m.handler.CreateCommandCenter)
	rg.GET("/disasters/:ref_id/command-center", m.handler.GetCommandCenterByDisaster)
}
