package disaster

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
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
	logger := logging.FromContext(c.Request.Context())

	var req CreateDisasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("create disaster: invalid request body", "error", err)
		response.Error(c, apperror.BadRequest("invalid_request_body", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}

	d, err := h.service.CreateDisaster(c.Request.Context(), req)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create disaster failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, toDisasterResponse(d))
}

func (h *Handler) ListDisasters(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	disasters, err := h.service.ListDisasters(c.Request.Context())
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("list disasters failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, toDisasterResponses(disasters))
}

func (h *Handler) GetDisaster(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("get disaster: invalid id", "id", c.Param("id"))
		response.Error(c, apperror.BadRequest("invalid_id", "disaster id must be a number"))
		return
	}

	d, err := h.service.GetDisaster(c.Request.Context(), int32(id))
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get disaster failed", "id", id, "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, toDisasterResponse(d))
}
