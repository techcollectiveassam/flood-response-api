package sos

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/pagination"
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

func (h *Handler) ListSOS(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	var query pagination.Query
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Warn("list sos: invalid query params", "error", err)
		response.Error(c, apperror.BadRequest("invalid_query_params", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}

	page, limit, err := pagination.Resolve(query, pagination.FromContext(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.service.ListSOS(c.Request.Context(), page, limit)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("list sos failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	items := make([]SOSResponse, 0, len(result.Items))
	for _, s := range result.Items {
		items = append(items, toSOSResponse(s))
	}

	paginationResp := response.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      result.Total,
		TotalPages: pagination.TotalPages(result.Total, limit),
	}
	response.Paginated(c, http.StatusOK, items, paginationResp)
}
