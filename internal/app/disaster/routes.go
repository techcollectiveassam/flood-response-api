package disaster

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/disasters", m.handler.CreateDisaster)
	rg.GET("/disasters", m.handler.ListDisasters)
}
