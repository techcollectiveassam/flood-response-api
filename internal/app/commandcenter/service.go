package commandcenter

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

var ErrInvalidDisasterRef = apperror.BadRequest("invalid_request_body", "disaster_ref_id must be a valid UUID")

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
		"disaster_ref_id", req.DisasterRefID,
		"name", req.Name,
		"type", req.Type,
	)

	refID, err := uuid.Parse(req.DisasterRefID)
	if err != nil {
		return nil, ErrInvalidDisasterRef
	}

	disasterID, err := s.repository.GetDisasterIDByRefID(ctx, refID)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("resolve disaster for command center", "error", err)
		}
		return nil, err
	}

	commandCenter := &CommandCenter{
		DisasterID:    disasterID,
		DisasterRefID: refID,
		Name:          req.Name,
		Type:          req.Type,
		Description:   req.Description,
		ContactPerson: req.ContactPerson,
		ContactMobile: req.ContactMobile,
		ContactEmail:  req.ContactEmail,
	}
	err = s.repository.Create(ctx, commandCenter)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create command center", "error", err)
		}
		return nil, err
	}
	return commandCenter, nil
}

func (s *Service) GetCommandCenterByDisaster(ctx context.Context, refID uuid.UUID) (*CommandCenter, error) {
	logger := logging.FromContext(ctx)

	// The disaster is the path resource, so resolve it before ever querying
	// the command center.
	disasterID, resolveErr := s.repository.GetDisasterIDByRefID(ctx, refID)
	if resolveErr != nil {
		if apperror.HTTPStatus(resolveErr) >= 500 {
			logger.Error("resolve disaster for command center", "ref_id", refID, "error", resolveErr)
		}
		return nil, resolveErr
	}

	commandCenter, err := s.repository.GetByDisasterID(ctx, disasterID)
	if err == nil {
		commandCenter.DisasterRefID = refID
		logger.Debug("command center found for disaster", "ref_id", refID, "id", commandCenter.ID)
		return commandCenter, nil
	}
	if !errors.Is(err, ErrCommandCenterNotFound) {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get command center by disaster", "ref_id", refID, "error", err)
		}
		return nil, err
	}

	logger.Warn("command center not found for disaster", "ref_id", refID)
	return nil, ErrCommandCenterNotFound
}
