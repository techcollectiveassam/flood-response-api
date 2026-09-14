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
	repository            Repository
	duplicateRadiusMeters float64
}

func NewService(repository Repository, duplicateRadiusMeters float64) *Service {
	return &Service{
		repository:            repository,
		duplicateRadiusMeters: duplicateRadiusMeters,
	}
}

// CreateSOS creates a new SOS request or merges into a nearby duplicate. The
// returned boolean reports whether a new row was created.
func (s *Service) CreateSOS(ctx context.Context, req CreateSOSRequest) (*SOS, bool, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("create sos request",
		"disaster_id", req.DisasterID,
		"latitude", req.Latitude,
		"longitude", req.Longitude,
		"reporter_mobile", req.ReporterMobile,
	)

	if err := validateCreateSOS(req); err != nil {
		return nil, false, err
	}

	if req.Latitude != nil {
		geometry := pointGeoJSON(req.Longitude, req.Latitude)
		matched, err := s.repository.FindNearby(ctx, req.DisasterID, req.ReporterMobile, geometry, s.duplicateRadiusMeters)
		if err != nil {
			if apperror.HTTPStatus(err) >= 500 {
				logger.Error("find nearby sos", "error", err)
			}
			return nil, false, err
		}

		if matched != nil {
			existing, err := s.repository.IncrementReportCount(ctx, matched.ID)
			if err != nil {
				if apperror.HTTPStatus(err) >= 500 {
					logger.Error("increment sos report count", "error", err)
				}
				return nil, false, err
			}

			logger.Debug("duplicate sos merged", "id", existing.ID, "report_count", existing.ReportCount)
			return existing, false, nil
		}
	}

	sosRequest := &SOS{
		DisasterID:     req.DisasterID,
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
		return nil, false, err
	}

	logger.Debug("sos created", "id", sosRequest.ID)
	return sosRequest, true, nil
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
