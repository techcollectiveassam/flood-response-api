package affectedarea

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

var ErrInvalidDisasterRef = apperror.BadRequest("invalid_request_body", "disaster_ref_id must be a valid UUID")

var severityRank = map[string]int{
	"unknown":  0,
	"low":      1,
	"medium":   2,
	"high":     3,
	"critical": 4,
}

type Service struct {
	repository Repository
	resolver   *AffectedAreaResolver
}

func NewService(repository Repository, resolver *AffectedAreaResolver) *Service {
	return &Service{
		repository: repository,
		resolver:   resolver,
	}
}

func (s *Service) CreateAffectedArea(ctx context.Context, req CreateAffectedAreaRequest) (*CreateAffectedAreaResponse, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("create affected area request",
		"disaster_ref_id", req.DisasterRefID,
		"severity", req.Severity,
	)

	if req.Location == nil {
		return nil, apperror.BadRequest("invalid_location", "location is required")
	}
	if err := req.Location.Validate(); err != nil {
		return nil, err
	}

	refID, err := uuid.Parse(req.DisasterRefID)
	if err != nil {
		return nil, ErrInvalidDisasterRef
	}

	disasterID, err := s.repository.GetDisasterIDByRefID(ctx, refID)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("resolve disaster for affected area", "error", err)
		}
		return nil, err
	}

	resolved, err := s.resolver.Resolve(ctx, disasterID, *req.Location)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("resolve affected area", "error", err)
		}
		return nil, err
	}

	var (
		area   *AffectedArea
		report *AffectedAreaReport
	)

	err = s.repository.WithTx(ctx, func(repository Repository) error {
		area, err = s.matchOrCreateArea(ctx, repository, req, resolved, disasterID, refID)
		if err != nil {
			if apperror.HTTPStatus(err) >= 500 {
				logger.Error("ensure affected area", "error", err)
			}
			return err
		}

		report = &AffectedAreaReport{
			AreaID:          area.ID,
			Description:     req.Description,
			LocationPayload: locationPayloadJSON(*req.Location),
			Severity:        req.Severity,
			ReporterName:    reporterName(req.Reporter),
			ReporterMobile:  reporterMobile(req.Reporter),
		}
		if err := repository.CreateReport(ctx, report); err != nil {
			if apperror.HTTPStatus(err) >= 500 {
				logger.Error("create affected area report", "error", err)
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if resolved.Decision == DecisionUncertain {
		logger.Warn("affected area match uncertain; needs investigation",
			"disaster_ref_id", req.DisasterRefID,
			"confidence", resolved.Confidence,
			"match_reason", resolved.MatchReason,
		)
	}

	result := &CreateAffectedAreaResponse{
		AffectedArea: toAffectedAreaResponse(area),
		Report:       toAffectedAreaReportResponse(report),
		IsNewArea:    resolved.Decision != DecisionMatched,
		Confidence:   resolved.Confidence,
		MatchReason:  resolved.MatchReason,
	}

	logger.Debug("affected area created",
		"area_id", area.ID,
		"is_new_area", result.IsNewArea,
		"confidence", result.Confidence,
		"match_reason", result.MatchReason,
	)

	return result, nil
}

func (s *Service) matchOrCreateArea(ctx context.Context, repository Repository, req CreateAffectedAreaRequest, resolved *ResolvedArea, disasterID int32, disasterRefID uuid.UUID) (*AffectedArea, error) {
	if resolved.Decision == DecisionMatched && resolved.Area != nil {
		if severityHigher(req.Severity, resolved.Area.Severity) {
			if err := repository.UpdateAreaSeverity(ctx, resolved.Area.ID, req.Severity); err != nil {
				return nil, err
			}
			resolved.Area.Severity = req.Severity
		}
		if err := repository.IncrementReportCount(ctx, resolved.Area.ID); err != nil {
			return nil, err
		}
		resolved.Area.ReportCount++
		return resolved.Area, nil
	}

	verificationStatus := StatusReported
	if resolved.Decision == DecisionUncertain {
		verificationStatus = StatusInvestigating
	}

	area := &AffectedArea{
		Description:        req.Description,
		DisasterID:         disasterID,
		DisasterRefID:      disasterRefID,
		Location:           locationLabel(*req.Location),
		Geometry:           locationGeometry(*req.Location),
		Latitude:           req.Location.Latitude,
		Longitude:          req.Location.Longitude,
		Severity:           req.Severity,
		VerificationStatus: verificationStatus,
	}

	if err := repository.CreateArea(ctx, area); err != nil {
		return nil, err
	}
	if err := repository.IncrementReportCount(ctx, area.ID); err != nil {
		return nil, err
	}
	area.ReportCount = 1
	return area, nil
}

func severityHigher(candidate, current string) bool {
	return severityRank[candidate] > severityRank[current]
}

type ListAffectedAreasResult struct {
	Areas []*AffectedArea
	Total int64
}

func (s *Service) ListAffectedAreas(ctx context.Context, page, limit int) (*ListAffectedAreasResult, error) {
	logger := logging.FromContext(ctx)
	areas, total, err := s.repository.ListAffectedAreas(ctx, page, limit)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			logger.Error("list affected areas", "error", err)
		}
		return nil, err
	}

	logger.Debug("affected areas listed", "count", len(areas), "total", total)
	return &ListAffectedAreasResult{Areas: areas, Total: total}, nil
}

func (s *Service) GetAffectedAreaByRefID(ctx context.Context, refID uuid.UUID) (*AffectedArea, error) {
	logger := logging.FromContext(ctx)

	area, err := s.repository.GetAffectedAreaByRefID(ctx, refID)
	if err != nil {
		if errors.Is(err, ErrAffectedAreaNotFound) {
			logger.Warn("affected area not found", "ref_id", refID)
		} else if apperror.HTTPStatus(err) >= 500 {
			logger.Error("get affected area", "ref_id", refID, "error", err)
		}
		return nil, err
	}

	logger.Debug("affected area found", "ref_id", refID)
	return area, nil
}
