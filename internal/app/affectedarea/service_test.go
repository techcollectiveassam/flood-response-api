package affectedarea

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testDisasterRefID = "11111111-1111-1111-1111-111111111111"

func newTestService(repo *mockRepository) *Service {
	return NewService(repo, NewAffectedAreaResolver())
}

func textRequest() CreateAffectedAreaRequest {
	return CreateAffectedAreaRequest{
		DisasterRefID: testDisasterRefID,
		Severity:      "medium",
		Location: &LocationPayload{
			Source:      SourceText,
			Description: "Village near the old bridge",
		},
	}
}

func floatPointer(value float64) *float64 {
	return &value
}

func TestCreateAffectedAreaTextSource(t *testing.T) {
	repo := &mockRepository{}
	svc := newTestService(repo)
	req := textRequest()

	result, err := svc.CreateAffectedArea(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, result.IsNewArea)
	assert.Equal(t, MatchReasonUnmatchable, result.MatchReason)
	assert.Equal(t, 0.0, result.Confidence)

	area := repo.createdArea
	assert.NotNil(t, area)
	assert.Equal(t, "Village near the old bridge", area.Location)
	assert.Equal(t, "", area.Geometry)
	assert.Equal(t, StatusReported, area.VerificationStatus)
	assert.Equal(t, int64(1), area.ReportCount)
	assert.Nil(t, area.Latitude)
	assert.Nil(t, area.Longitude)
	assert.Equal(t, int32(1), area.DisasterID)
	assert.Equal(t, testDisasterRefID, area.DisasterRefID.String())
	assert.Equal(t, 1, repo.resolveCalls)
	assert.Equal(t, testDisasterRefID, repo.resolvedRefID.String())

	report := repo.createdReport
	assert.NotNil(t, report)
	assert.Equal(t, area.ID, report.AreaID)
	assert.NotEqual(t, uuid.Nil, report.RefID)
	assert.Equal(t, "medium", report.Severity)
	assert.Empty(t, report.ReporterName)
	assert.Empty(t, report.ReporterMobile)
}

func TestCreateAffectedAreaUnknownDisaster(t *testing.T) {
	repo := &mockRepository{resolveErr: ErrDisasterNotFound}
	svc := newTestService(repo)

	result, err := svc.CreateAffectedArea(context.Background(), textRequest())

	assert.ErrorIs(t, err, ErrDisasterNotFound)
	assert.Nil(t, result)
	assert.Equal(t, 1, repo.resolveCalls)
	assert.Nil(t, repo.createdArea)
}

func TestCreateAffectedAreaInvalidDisasterRefID(t *testing.T) {
	repo := &mockRepository{}
	svc := newTestService(repo)

	req := textRequest()
	req.DisasterRefID = "not-a-uuid"

	result, err := svc.CreateAffectedArea(context.Background(), req)

	assert.ErrorIs(t, err, ErrInvalidDisasterRef)
	assert.Nil(t, result)
	assert.Equal(t, 0, repo.resolveCalls)
	assert.Nil(t, repo.createdArea)
}

func TestUpdateAreaSeverity(t *testing.T) {
	tests := []struct {
		name      string
		candidate string
		current   string
		want      bool
	}{
		{name: "higher", candidate: "critical", current: "medium", want: true},
		{name: "equal", candidate: "medium", current: "medium", want: false},
		{name: "lower", candidate: "low", current: "high", want: false},
		{name: "low beats unknown", candidate: "low", current: "unknown", want: true},
		{name: "unknown does not beat low", candidate: "unknown", current: "low", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, severityHigher(tt.candidate, tt.current))
		})
	}
}

func TestListAffectedAreas(t *testing.T) {
	areas := []*AffectedArea{
		{ID: 2, Severity: "high", VerificationStatus: StatusReported, ReportCount: 1},
		{ID: 1, Severity: "medium", VerificationStatus: StatusReported, ReportCount: 5},
	}

	repo := &mockRepository{listAreas: areas, listTotal: 42}
	svc := newTestService(repo)

	result, err := svc.ListAffectedAreas(context.Background(), 2, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), result.Total)
	assert.Len(t, result.Areas, 2)
	assert.Equal(t, 2, repo.listPage)
	assert.Equal(t, 10, repo.listLimit)
}

func TestListAffectedAreasError(t *testing.T) {
	repo := &mockRepository{listErr: assert.AnError}
	svc := newTestService(repo)

	result, err := svc.ListAffectedAreas(context.Background(), 1, 20)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAffectedArea(t *testing.T) {
	refID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	area := &AffectedArea{ID: 1, RefID: refID, Severity: "high", VerificationStatus: StatusReported, ReportCount: 3}
	repo := &mockRepository{getAffected: area}
	svc := newTestService(repo)

	result, err := svc.GetAffectedAreaByRefID(context.Background(), refID)

	assert.NoError(t, err)
	assert.Equal(t, area, result)
}

func TestGetAffectedAreaNotFound(t *testing.T) {
	svc := newTestService(&mockRepository{})

	result, err := svc.GetAffectedAreaByRefID(context.Background(), uuid.MustParse("99999999-9999-9999-9999-999999999999"))

	assert.ErrorIs(t, err, ErrAffectedAreaNotFound)
	assert.Nil(t, result)
}

func TestGetAffectedAreaError(t *testing.T) {
	repo := &mockRepository{getErr: assert.AnError}
	svc := newTestService(repo)

	result, err := svc.GetAffectedAreaByRefID(context.Background(), uuid.MustParse("55555555-5555-5555-5555-555555555555"))

	assert.Error(t, err)
	assert.Nil(t, result)
}
