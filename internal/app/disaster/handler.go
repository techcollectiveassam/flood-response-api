package disaster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	refID, err := uuid.Parse(c.Param("ref_id"))
	if err != nil {
		logger.Warn("get disaster: invalid ref_id", "ref_id", c.Param("ref_id"))
		response.Error(c, apperror.BadRequest("invalid_id", "disaster ref_id must be a valid UUID"))
		return
	}

	d, err := h.service.GetDisaster(c.Request.Context(), refID)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get disaster failed", "ref_id", refID, "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, toDisasterResponse(d))
}
