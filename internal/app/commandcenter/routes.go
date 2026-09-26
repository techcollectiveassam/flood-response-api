package commandcenter

import (
	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/command-centers", m.handler.CreateCommandCenter)
}
