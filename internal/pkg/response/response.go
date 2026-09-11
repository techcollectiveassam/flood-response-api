package response

import (
	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

func Data(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func Error(c *gin.Context, err error) {
	body := gin.H{
		"code":    apperror.Code(err),
		"message": apperror.Message(err),
	}
	if details := apperror.Details(err); len(details) > 0 {
		body["details"] = details
	}
	c.JSON(apperror.HTTPStatus(err), gin.H{"error": body})
}
