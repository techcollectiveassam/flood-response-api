package sos

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/sos", m.handler.CreateSOS)
}
