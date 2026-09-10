package disaster

import (
	"net/http"
	"strconv"

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
		response.Error(c, apperror.BadRequest("invalid_request_body", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}

	d, err := h.service.CreateDisaster(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, toDisasterResponse(d))
}

func (h *Handler) ListDisasters(c *gin.Context) {
	disasters, err := h.service.ListDisasters(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, toDisasterResponses(disasters))
}

func (h *Handler) GetDisaster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, apperror.BadRequest("invalid_id", "disaster id must be a number"))
		return
	}

	d, err := h.service.GetDisaster(c.Request.Context(), int32(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, toDisasterResponse(d))
}
