package sos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSOS(t *testing.T) {
	tests := []struct {
		name           string
		req            CreateSOSRequest
		findNearby     *SOS
		findNearbyErr  error
		repoErr        error
		wantErr        error
		wantCreated    bool
		wantFindNearby bool
		wantMobile     string
		wantIncrements int
		wantCreates    int
	}{
		{
			name: "creates when no nearby match",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			wantFindNearby: true,
			wantCreated:    true,
			wantCreates:    1,
		},
		{
			name: "merges into nearby match",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			findNearby:     &SOS{ID: 9, DisasterID: 7, ReportCount: 3},
			wantFindNearby: true,
			wantCreated:    false,
			wantIncrements: 1,
			wantCreates:    0,
		},
		{
			name: "mobile scopes the nearby lookup",
			req: CreateSOSRequest{
				DisasterID:     7,
				ReporterMobile: "9876543210",
				Location:       location(26.18, 91.73),
			},
			findNearby:     &SOS{ID: 9, DisasterID: 7, ReportCount: 3},
			wantFindNearby: true,
			wantMobile:     "9876543210",
			wantCreated:    false,
			wantIncrements: 1,
			wantCreates:    0,
		},
		{
			name: "mobile only skips nearby lookup",
			req: CreateSOSRequest{
				DisasterID:     7,
				ReporterMobile: "9876543210",
			},
			wantFindNearby: false,
			wantCreated:    true,
			wantCreates:    1,
		},
		{
			name: "empty location with mobile creates without geometry",
			req: CreateSOSRequest{
				DisasterID:     7,
				ReporterMobile: "9876543210",
				Location:       &Location{},
			},
			wantFindNearby: false,
			wantCreated:    true,
			wantCreates:    1,
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
			name: "find nearby error",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			findNearbyErr:  assert.AnError,
			wantErr:        assert.AnError,
			wantFindNearby: true,
		},
		{
			name: "create repository error",
			req: CreateSOSRequest{
				DisasterID: 7,
				Location:   location(26.18, 91.73),
			},
			repoErr:        assert.AnError,
			wantErr:        assert.AnError,
			wantFindNearby: true,
			wantCreates:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				findNearbyResult: tt.findNearby,
				findNearbyErr:    tt.findNearbyErr,
				createErr:        tt.repoErr,
			}
			svc := NewService(repo, 100)

			result, created, err := svc.CreateSOS(context.Background(), tt.req)

			assert.Equal(t, tt.wantFindNearby, repo.findNearbyCalled)
			assert.Equal(t, tt.wantMobile, repo.findNearbyMobile)
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
