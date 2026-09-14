package sos

import (
	"context"
	"time"
)

type mockRepository struct {
	findNearbyResult *SOS
	findNearbyErr    error
	incrementErr     error
	createErr        error

	findNearbyCalled bool
	findNearbyMobile string
	incrementCalls   int
	createCalls      int
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

func (m *mockRepository) FindNearby(ctx context.Context, disasterID int32, mobile, geometry string, radiusMeters float64) (*SOS, error) {
	m.findNearbyCalled = true
	m.findNearbyMobile = mobile
	if m.findNearbyErr != nil {
		return nil, m.findNearbyErr
	}
	return m.findNearbyResult, nil
}

func (m *mockRepository) IncrementReportCount(ctx context.Context, id int64) (*SOS, error) {
	m.incrementCalls++
	if m.incrementErr != nil {
		return nil, m.incrementErr
	}
	s := *m.findNearbyResult
	s.ReportCount++
	return &s, nil
}
