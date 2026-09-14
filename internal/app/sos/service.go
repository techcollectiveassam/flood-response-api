package sos

import (
	"context"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

var (
	ErrInvalidSOS        = apperror.BadRequest("invalid_sos", "either reporter_mobile or location is required")
	ErrInvalidLocation   = apperror.BadRequest("invalid_location", "latitude and longitude must be provided together")
	ErrInvalidCoordinate = apperror.BadRequest("invalid_location", "latitude must be within -90..90 and longitude within -180..180")
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateSOS(ctx context.Context, req CreateSOSRequest) (*SOS, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("create sos request",
		"latitude", req.Latitude,
		"longitude", req.Longitude,
		"reporter_mobile", req.ReporterMobile,
	)

	if err := validateCreateSOS(req); err != nil {
		return nil, err
	}

	sosRequest := &SOS{
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		ReporterMobile: req.ReporterMobile,
		Message:        req.Message,
		Status:         StatusReported,
	}

	if err := s.repository.Create(ctx, sosRequest); err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("create sos", "error", err)
		}
		return nil, err
	}

	logger.Debug("sos created", "id", sosRequest.ID)
	return sosRequest, nil
}

func validateCreateSOS(req CreateSOSRequest) error {
	hasMobile := req.ReporterMobile != ""
	hasLatitude := req.Latitude != nil
	hasLongitude := req.Longitude != nil

	if !hasMobile && !hasLatitude && !hasLongitude {
		return ErrInvalidSOS
	}
	if hasLatitude != hasLongitude {
		return ErrInvalidLocation
	}
	if hasLatitude {
		if *req.Latitude < -90 || *req.Latitude > 90 || *req.Longitude < -180 || *req.Longitude > 180 {
			return ErrInvalidCoordinate
		}
	}
	return nil
}
