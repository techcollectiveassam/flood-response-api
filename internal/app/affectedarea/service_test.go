package affectedarea

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestService(repo *mockRepository) *Service {
	return NewService(repo, NewAffectedAreaResolver())
}

func textRequest() CreateAffectedAreaRequest {
	return CreateAffectedAreaRequest{
		DisasterID: 1,
		Name:       "Possible flooding",
		Severity:   "medium",
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

	report := repo.createdReport
	assert.NotNil(t, report)
	assert.Equal(t, area.ID, report.AreaID)
	assert.Equal(t, "medium", report.Severity)
	assert.Empty(t, report.ReporterName)
	assert.Empty(t, report.ReporterMobile)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, severityHigher(tt.candidate, tt.current))
		})
	}
}
