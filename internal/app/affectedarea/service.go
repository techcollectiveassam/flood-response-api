package affectedarea

import (
	"context"
	"log/slog"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
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
	logger     *slog.Logger
}

func NewService(repository Repository, resolver *AffectedAreaResolver, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		resolver:   resolver,
		logger:     logger,
	}
}

func (s *Service) CreateAffectedArea(ctx context.Context, req CreateAffectedAreaRequest) (*CreateAffectedAreaResponse, error) {
	if err := req.Location.Validate(); err != nil {
		return nil, err
	}

	resolved, err := s.resolver.Resolve(ctx, req.DisasterID, req.Location)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			s.logger.Error("resolve affected area", "error", err)
		}
		return nil, err
	}

	area, err := s.ensureArea(ctx, req, resolved)
	if err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			s.logger.Error("ensure affected area", "error", err)
		}
		return nil, err
	}

	report := &AffectedAreaReport{
		AreaID:          area.ID,
		Name:            req.Name,
		Description:     req.Description,
		LocationPayload: locationPayloadJSON(req.Location),
		Severity:        req.Severity,
		ReporterName:    reporterName(req.Reporter),
		ReporterMobile:  reporterMobile(req.Reporter),
	}
	if err := s.repository.CreateReport(ctx, report); err != nil {
		if apperror.HTTPStatus(err) >= 500 {
			s.logger.Error("create affected area report", "error", err)
		}
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

func (s *Service) ensureArea(ctx context.Context, req CreateAffectedAreaRequest, resolved *ResolvedArea) (*AffectedArea, error) {
	if resolved.Decision == DecisionMatched && resolved.Area != nil {
		if severityHigher(req.Severity, resolved.Area.Severity) {
			if err := s.repository.UpdateAreaSeverity(ctx, resolved.Area.ID, req.Severity); err != nil {
				return nil, err
			}
			resolved.Area.Severity = req.Severity
		}
		if err := s.repository.IncrementReportCount(ctx, resolved.Area.ID); err != nil {
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
		Location:           locationLabel(req.Location),
		Geometry:           locationGeometry(req.Location),
		Latitude:           req.Location.Latitude,
		Longitude:          req.Location.Longitude,
		Severity:           req.Severity,
		VerificationStatus: verificationStatus,
	}

	if err := s.repository.CreateArea(ctx, area); err != nil {
		return nil, err
	}
	if err := s.repository.IncrementReportCount(ctx, area.ID); err != nil {
		return nil, err
	}
	area.ReportCount = 1
	return area, nil
}

func severityHigher(candidate, current string) bool {
	return severityRank[candidate] > severityRank[current]
}
