package affectedarea

import (
	"context"
	"log/slog"
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

func (s *Service) CreateAffectedArea(ctx context.Context, req CreateAffectedAreaRequest) error {
	area := &AffectedArea{
		Name:        req.Name,
		Description: req.Description,
		DisasterID:  req.DisasterID,
		Location:    req.Location,
		Severity:    req.Severity,
		Source:      req.Source,
	}

	// Validation and other domain rules belong here, but are outside this skeleton.
	return s.repository.Create(ctx, area)
}
