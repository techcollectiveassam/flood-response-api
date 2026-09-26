package commandcenter

import (
	"context"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateCommandCenter(ctx context.Context, req CreateCommandCenterRequest) (*CommandCenter, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("create command center request",
		"name", req.Name,
		"type", req.Type,
	)

	commandCenter := &CommandCenter{
		DisasterID:    req.DisasterID,
		Name:          req.Name,
		Type:          req.Type,
		Description:   req.Description,
		ContactPerson: req.ContactPerson,
		ContactMobile: req.ContactMobile,
		ContactEmail:  req.ContactEmail,
	}
	err := s.repository.Create(ctx, commandCenter)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create command center", "error", err)
		}
		return nil, err
	}
	return commandCenter, nil
}
