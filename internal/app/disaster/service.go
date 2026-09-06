package disaster

import (
	"context"
	"log/slog"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
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
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
	}

	if err := s.repository.Create(ctx, d); err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			s.logger.Error("create disaster", "error", err)
		}
		return nil, err
	}

	return d, nil
}
