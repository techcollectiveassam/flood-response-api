package sos

import (
	"net/http"

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

func (h *Handler) CreateSOS(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	var req CreateSOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("create sos: invalid request body", "error", err)
		response.Error(c, apperror.BadRequest("invalid_request_body", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}

	result, created, err := h.service.CreateSOS(c.Request.Context(), req)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create sos failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	response.Data(c, status, toSOSResponse(result))
}
