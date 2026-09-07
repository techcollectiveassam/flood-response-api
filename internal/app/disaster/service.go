package disaster

import (
	"context"
	"log/slog"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/timeutil"
)

type Service struct {
	repository Repository
	logger     *slog.Logger
}

func NewService(repository Repository, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) CreateDisaster(ctx context.Context, req CreateDisasterRequest) (*Disaster, error) {
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
			s.logger.Error("create disaster", "error", err)
		}
		return nil, err
	}

	return d, nil
}

func (s *Service) ListDisasters(ctx context.Context) ([]Disaster, error) {
	disasters, err := s.repository.List(ctx)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			s.logger.Error("list disasters", "error", err)
		}
		return nil, err
	}
	return disasters, nil
}
