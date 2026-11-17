package sos

import (
	"context"
	"time"
)

type mockRepository struct {
	findDuplicateResult *SOS
	findDuplicateErr    error
	incrementErr        error
	createErr           error
	listResults         []*SOS
	listTotal           int64
	listErr             error

	findDuplicateCalled bool
	findDuplicateMobile string
	findDuplicateGeom   string
	findDuplicateRadius float64
	incrementCalls      int
	createCalls         int
}

func (m *mockRepository) Create(ctx context.Context, sos *SOS) error {
	m.createCalls++
	if m.createErr != nil {
		return m.createErr
	}
	sos.ID = 1
	sos.DisasterID = 7
	sos.Status = StatusReported
	sos.CreatedAt = time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	return nil
}

func (m *mockRepository) FindDuplicate(ctx context.Context, disasterID int32, mobile, geometry string, radiusMeters float64) (*SOS, error) {
	m.findDuplicateCalled = true
	m.findDuplicateMobile = mobile
	m.findDuplicateGeom = geometry
	m.findDuplicateRadius = radiusMeters
	if m.findDuplicateErr != nil {
		return nil, m.findDuplicateErr
	}
	return m.findDuplicateResult, nil
}

func (m *mockRepository) IncrementReportCount(ctx context.Context, id int64) (*SOS, error) {
	m.incrementCalls++
	if m.incrementErr != nil {
		return nil, m.incrementErr
	}
	s := *m.findDuplicateResult
	s.ReportCount++
	return &s, nil
}

func (m *mockRepository) ListSOS(ctx context.Context, page, limit int) ([]*SOS, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.listResults, m.listTotal, nil
}
