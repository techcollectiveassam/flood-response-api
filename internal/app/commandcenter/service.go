package commandcenter

import (
	"context"
	"errors"

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

func (s *Service) GetCommandCenterByDisaster(ctx context.Context, disasterID int32) (*CommandCenter, error) {
	logger := logging.FromContext(ctx)

	// The disaster is the path resource, so validate the parent before ever
	// querying the command center.
	exists, existsErr := s.repository.DisasterExists(ctx, disasterID)
	if existsErr != nil {
		if apperror.HTTPStatus(existsErr) >= 500 {
			logger.Error("check disaster for command center", "disaster_id", disasterID, "error", existsErr)
		}
		return nil, existsErr
	}
	if !exists {
		logger.Warn("disaster not found", "disaster_id", disasterID)
		return nil, ErrDisasterNotFound
	}

	commandCenter, err := s.repository.GetByDisasterID(ctx, disasterID)
	if err == nil {
		logger.Debug("command center found for disaster", "disaster_id", disasterID, "id", commandCenter.ID)
		return commandCenter, nil
	}
	if !errors.Is(err, ErrCommandCenterNotFound) {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get command center by disaster", "disaster_id", disasterID, "error", err)
		}
		return nil, err
	}

	logger.Warn("command center not found for disaster", "disaster_id", disasterID)
	return nil, ErrCommandCenterNotFound
}
