package affectedarea

import "context"

type mockRepository struct {
	createAreaErr     error
	createReportErr   error
	updateSeverityErr error
	incrementErr      error
	listErr           error
	getErr            error

	createdArea   *AffectedArea
	createdReport *AffectedAreaReport

	listAreas   []*AffectedArea
	listTotal   int64
	getAffected *AffectedArea

	severityUpdates []string
	increments      []int64
	listPage        int
	listLimit       int
}

func (m *mockRepository) WithTx(ctx context.Context, fn func(Repository) error) error {
	return fn(m)
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

func (m *mockRepository) GetAffectedArea(ctx context.Context, id int64) (*AffectedArea, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.getAffected == nil {
		return nil, ErrAffectedAreaNotFound
	}
	return m.getAffected, nil
}
