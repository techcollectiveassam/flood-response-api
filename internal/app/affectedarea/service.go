package affectedarea

import (
	"context"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

var severityRank = map[string]int{
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

	if req.Location == nil {
		return nil, apperror.BadRequest("invalid_location", "location is required")
	}
	if err := req.Location.Validate(); err != nil {
		return nil, err
	}

	resolved, err := s.resolver.Resolve(ctx, req.DisasterID, *req.Location)
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
		area, err = s.matchOrCreateArea(ctx, repository, req, resolved)
		if err != nil {
			if apperror.HTTPStatus(err) >= 500 {
				logger.Error("ensure affected area", "error", err)
			}
			return err
		}

		report = &AffectedAreaReport{
			AreaID:          area.ID,
			Name:            req.Name,
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

	return &CreateAffectedAreaResponse{
		AffectedArea: toAffectedAreaResponse(area),
		Report:       toAffectedAreaReportResponse(report),
		IsNewArea:    resolved.Decision != DecisionMatched,
		Confidence:   resolved.Confidence,
		MatchReason:  resolved.MatchReason,
	}, nil
}

func (s *Service) matchOrCreateArea(ctx context.Context, repository Repository, req CreateAffectedAreaRequest, resolved *ResolvedArea) (*AffectedArea, error) {
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
		Name:               req.Name,
		Description:        req.Description,
		DisasterID:         req.DisasterID,
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
