package response

import (
	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

func Error(c *gin.Context, err error) {
	c.JSON(apperror.HTTPStatus(err), gin.H{
		"error": gin.H{
			"code":    apperror.Code(err),
			"message": apperror.Message(err),
		},
	})
}