package commandcenter

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

func (h *Handler) CreateCommandCenter(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	var req CreateCommandCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("create command center: invalid request body", "error", err)
		response.Error(c, apperror.BadRequest("invalid_request_body", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}
	result, err := h.service.CreateCommandCenter(c.Request.Context(), req)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create command center failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, toCommandCenterResponse(result))
}

func (h *Handler) GetCommandCenterByDisaster(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	refID, err := uuid.Parse(c.Param("ref_id"))
	if err != nil {
		logger.Warn("get command center by disaster: invalid ref_id", "ref_id", c.Param("ref_id"))
		response.Error(c, apperror.BadRequest("invalid_id", "disaster ref_id must be a valid UUID"))
		return
	}

	result, err := h.service.GetCommandCenterByDisaster(c.Request.Context(), refID)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get command center by disaster failed", "ref_id", refID, "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, toCommandCenterResponse(result))
}
