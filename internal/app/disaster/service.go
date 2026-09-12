package disaster

import (
	"context"
	"errors"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/timeutil"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateDisaster(ctx context.Context, req CreateDisasterRequest) (*Disaster, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("create disaster request",
		"name", req.Name,
		"type", req.Type,
		"starts_at", req.StartsAt,
		"ends_at", req.EndsAt,
	)

	d := &Disaster{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      StatusActive,
		StartsAt:    timeutil.UTC(req.StartsAt),
		EndsAt:      timeutil.UTC(req.EndsAt),
	}

	if err := s.repository.Create(ctx, d); err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create disaster", "error", err)
		}
		return nil, err
	}

	logger.Debug("disaster created", "id", d.ID)
	return d, nil
}

func (s *Service) ListDisasters(ctx context.Context) ([]Disaster, error) {
	logger := logging.FromContext(ctx)

	disasters, err := s.repository.List(ctx)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("list disasters", "error", err)
		}
		return nil, err
	}

	logger.Debug("disasters listed", "count", len(disasters))
	return disasters, nil
}

func (s *Service) GetDisaster(ctx context.Context, id int32) (*Disaster, error) {
	logger := logging.FromContext(ctx)

	d, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDisasterNotFound) {
			logger.Warn("disaster not found", "id", id)
		} else if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get disaster", "id", id, "error", err)
		}
		return nil, err
	}

	logger.Debug("disaster found", "id", id)
	return d, nil
}
