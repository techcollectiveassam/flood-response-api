package affectedarea

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

func (s *Service) CreateAffectedArea(ctx context.Context, req CreateAffectedAreaRequest) (*AffectedAreaResponse, error) {
	area := &AffectedArea{
		Name:        req.Name,
		Description: req.Description,
		DisasterID:  req.DisasterID,
		Location:    req.Location,
		Geometry:    req.Geometry,
		Severity:    req.Severity,
	}

	if err := s.repository.Create(ctx, area); err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			s.logger.Error("create affected area", "error", err)
		}
		return nil, err
	}

	return &AffectedAreaResponse{
		ID:          area.ID,
		Name:        area.Name,
		Description: area.Description,
		DisasterID:  area.DisasterID,
		Location:    area.Location,
		Geometry:    area.Geometry,
		Severity:    area.Severity,
	}, nil
}
