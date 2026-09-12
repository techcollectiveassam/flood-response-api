package affectedarea

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

func (h *Handler) CreateAffectedArea(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	var req CreateAffectedAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("create affected area: invalid request body", "error", err)
		response.Error(c, apperror.BadRequest("invalid_request_body", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}

	result, err := h.service.CreateAffectedArea(c.Request.Context(), req)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create affected area failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, result)
}

func (h *Handler) ListAffectedAreas(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	var query pagination.Query
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Warn("list affected areas: invalid query params", "error", err)
		response.Error(c, apperror.BadRequest("invalid_query_params", validation.MessageValidationFailed).WithDetails(validation.Details(err)))
		return
	}

	page, limit, err := pagination.Resolve(query, pagination.FromContext(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.service.ListAffectedAreas(c.Request.Context(), page, limit)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("list affected areas failed", "error", err)
		}
		response.Error(c, err)
		return
	}

	items := make([]AffectedAreaResponse, 0, len(result.Areas))
	for _, area := range result.Areas {
		items = append(items, toAffectedAreaResponse(area))
	}

	paginationResp := response.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      result.Total,
		TotalPages: pagination.TotalPages(result.Total, limit),
	}
	response.Paginated(c, http.StatusOK, items, paginationResp)
}
