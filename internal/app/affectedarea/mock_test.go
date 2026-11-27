package affectedarea

import (
	"context"

	"github.com/google/uuid"
)

type mockRepository struct {
	createAreaErr     error
	createReportErr   error
	updateSeverityErr error
	incrementErr      error
	listErr           error
	getErr            error
	resolveErr        error

	disasterID int32

	createdArea   *AffectedArea
	createdReport *AffectedAreaReport

	listAreas   []*AffectedArea
	listTotal   int64
	getAffected *AffectedArea

	severityUpdates []string
	increments      []int64
	listPage        int
	listLimit       int
	resolveCalls    int
	resolvedRefID   uuid.UUID
}

func (m *mockRepository) WithTx(ctx context.Context, fn func(Repository) error) error {
	return fn(m)
}

func (m *mockRepository) GetDisasterIDByRefID(ctx context.Context, refID uuid.UUID) (int32, error) {
	m.resolveCalls++
	m.resolvedRefID = refID
	if m.resolveErr != nil {
		return 0, m.resolveErr
	}
	if m.disasterID == 0 {
		return 1, nil
	}
	return m.disasterID, nil
}

func (m *mockRepository) CreateArea(ctx context.Context, area *AffectedArea) error {
	if m.createAreaErr != nil {
		return m.createAreaErr
	}
	if area.ID == 0 {
		area.ID = 1
	}
	m.createdArea = area
	return nil
}

func (m *mockRepository) CreateReport(ctx context.Context, report *AffectedAreaReport) error {
	if m.createReportErr != nil {
		return m.createReportErr
	}
	if report.ID == 0 {
		report.ID = 1
	}
	if report.RefID == uuid.Nil {
		report.RefID = uuid.MustParse("66666666-6666-6666-6666-666666666666")
	}
	m.createdReport = report
	return nil
}

func (m *mockRepository) UpdateAreaSeverity(ctx context.Context, id int64, severity string) error {
	if m.updateSeverityErr != nil {
		return m.updateSeverityErr
	}
	m.severityUpdates = append(m.severityUpdates, severity)
	return nil
}

func (m *mockRepository) IncrementReportCount(ctx context.Context, id int64) error {
	if m.incrementErr != nil {
		return m.incrementErr
	}
	m.increments = append(m.increments, id)
	return nil
}

func (m *mockRepository) ListAffectedAreas(ctx context.Context, page, limit int) ([]*AffectedArea, int64, error) {
	m.listPage = page
	m.listLimit = limit
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.listAreas, m.listTotal, nil
}

func (m *mockRepository) GetAffectedAreaByRefID(ctx context.Context, refID uuid.UUID) (*AffectedArea, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.getAffected == nil {
		return nil, ErrAffectedAreaNotFound
	}
	return m.getAffected, nil
}
