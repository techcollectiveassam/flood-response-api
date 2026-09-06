package disaster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/response"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/validation"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateDisaster(c *gin.Context) {
	var req CreateDisasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("invalid_request_body", validation.Message(err)))
		return
	}

	d, err := h.service.CreateDisaster(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, toDisasterResponse(d))
}
