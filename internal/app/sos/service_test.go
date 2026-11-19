package sos

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSOS(t *testing.T) {
	tests := []struct {
		name              string
		req               CreateSOSRequest
		findDuplicate     *SOS
		findDuplicateErr  error
		repoErr           error
		wantErr           error
		wantCreated       bool
		wantFindDuplicate bool
		wantMobile        string
		wantGeometry      string
		wantRadius        float64
		checkQueryArgs    bool
		wantIncrements    int
		wantCreates       int
	}{
		{
			name: "creates when no nearby match",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			wantFindDuplicate: true,
			wantCreated:       true,
			wantCreates:       1,
		},
		{
			name: "merges into nearby match",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			findDuplicate:     &SOS{ID: 9, DisasterID: 7, ReportCount: 3},
			wantFindDuplicate: true,
			wantCreated:       false,
			wantIncrements:    1,
			wantCreates:       0,
		},
		{
			name: "same mobile matches",
			req: CreateSOSRequest{
				DisasterID:     7,
				ReporterMobile: "0000000000",
				Location:       location(26.18, 91.73),
			},
			findDuplicate:     &SOS{ID: 9, DisasterID: 7, ReportCount: 3},
			wantFindDuplicate: true,
			wantMobile:        "0000000000",
			wantGeometry:      `{"type":"Point","coordinates":[91.73,26.18]}`,
			wantRadius:        100,
			checkQueryArgs:    true,
			wantCreated:       false,
			wantIncrements:    1,
			wantCreates:       0,
		},
		{
			name: "mobile only checks for duplicate",
			req: CreateSOSRequest{
				DisasterID:     7,
				ReporterMobile: "0000000000",
			},
			findDuplicate:     &SOS{ID: 9, DisasterID: 7, ReportCount: 3},
			wantFindDuplicate: true,
			wantMobile:        "0000000000",
			wantRadius:        100,
			checkQueryArgs:    true,
			wantCreated:       false,
			wantIncrements:    1,
			wantCreates:       0,
		},
		{
			name: "empty location with mobile creates without geometry",
			req: CreateSOSRequest{
				DisasterID:     7,
				ReporterMobile: "0000000000",
				Location:       &Location{},
			},
			wantFindDuplicate: true,
			wantMobile:        "0000000000",
			wantCreated:       true,
			wantCreates:       1,
		},
		{
			name: "missing mobile and location",
			req: CreateSOSRequest{
				DisasterID: 7,
			},
			wantErr: ErrInvalidSOS,
		},
		{
			name: "empty location object without mobile",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   &Location{},
			},
			wantErr: ErrInvalidSOS,
		},
		{
			name: "partial location",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   &Location{Latitude: fp(26.18)},
			},
			wantErr: ErrInvalidLocation,
		},
		{
			name: "out of range coordinate",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(95, 91.73),
			},
			wantErr: ErrInvalidCoordinate,
		},
		{
			name: "find duplicate error",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			findDuplicateErr:  assert.AnError,
			wantErr:           assert.AnError,
			wantFindDuplicate: true,
		},
		{
			name: "create repository error",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			repoErr:           assert.AnError,
			wantErr:           assert.AnError,
			wantFindDuplicate: true,
			wantCreates:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				findDuplicateResult: tt.findDuplicate,
				findDuplicateErr:    tt.findDuplicateErr,
				createErr:           tt.repoErr,
			}
			svc := NewService(repo, 100)

			result, created, err := svc.CreateSOS(context.Background(), tt.req)

			assert.Equal(t, tt.wantFindDuplicate, repo.findDuplicateCalled)
			assert.Equal(t, tt.wantMobile, repo.findDuplicateMobile)
			if tt.checkQueryArgs {
				assert.Equal(t, tt.wantGeometry, repo.findDuplicateGeom)
				assert.Equal(t, tt.wantRadius, repo.findDuplicateRadius)
			}
			assert.Equal(t, tt.wantIncrements, repo.incrementCalls)
			assert.Equal(t, tt.wantCreates, repo.createCalls)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
				assert.False(t, created)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCreated, created)
			if tt.wantCreated {
				assert.Equal(t, tt.req.DisasterID, result.DisasterID)
			} else {
				assert.Equal(t, int32(4), result.ReportCount)
			}
		})
	}
}

func location(latitude, longitude float64) *Location {
	return &Location{Latitude: fp(latitude), Longitude: fp(longitude)}
}

func fp(value float64) *float64 {
	return &value
}

func TestGetSOSByTrackingID(t *testing.T) {
	trackingID := uuid.MustParse("3f2a9c1e-7b4d-4e2a-9c1e-5d6f7a8b9c0d")
	found := &SOS{ID: 9, TrackingID: trackingID, DisasterID: 7, Status: StatusReported, ReportCount: 2}

	tests := []struct {
		name    string
		result  *SOS
		repoErr error
		wantErr error
	}{
		{
			name:   "returns the sos for a known tracking id",
			result: found,
		},
		{
			name:    "returns not found when no sos matches",
			result:  nil,
			wantErr: ErrSOSNotFound,
		},
		{
			name:    "propagates repository errors",
			repoErr: assert.AnError,
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				getByTrackingResult: tt.result,
				getByTrackingErr:    tt.repoErr,
			}
			svc := NewService(repo, 100)

			result, err := svc.GetSOSByTrackingID(context.Background(), trackingID)

			assert.Equal(t, 1, repo.getByTrackingCalls)
			assert.Equal(t, trackingID, repo.getByTrackingID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, found, result)
		})
	}
}

func TestListSOS(t *testing.T) {
	tests := []struct {
		name      string
		results   []*SOS
		total     int64
		repoErr   error
		wantErr   error
		wantTotal int64
		wantLen   int
	}{
		{
			name: "success",
			results: []*SOS{
				{ID: 2, DisasterID: 7, Status: "reported", ReportCount: 2},
				{ID: 1, DisasterID: 8, Status: "reported", ReportCount: 1},
			},
			total:     2,
			wantTotal: 2,
			wantLen:   2,
		},
		{
			name:      "empty",
			results:   []*SOS{},
			total:     0,
			wantTotal: 0,
			wantLen:   0,
		},
		{
			name:    "repository error",
			repoErr: assert.AnError,
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{listResults: tt.results, listTotal: tt.total, listErr: tt.repoErr}
			svc := NewService(repo, 100)

			result, err := svc.ListSOS(context.Background(), 1, 20)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTotal, result.Total)
			assert.Len(t, result.Items, tt.wantLen)
		})
	}
}
