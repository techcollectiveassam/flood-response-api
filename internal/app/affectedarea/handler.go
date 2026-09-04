package affectedarea

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAffectedArea(c *gin.Context) {
	var req CreateAffectedAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CreateAffectedArea(c.Request.Context(), req); err != nil {
		if errors.Is(err, ErrNotImplemented) {
			response.Error(c, http.StatusNotImplemented, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "unable to create affected area")
		return
	}

	c.Status(http.StatusCreated)
}
