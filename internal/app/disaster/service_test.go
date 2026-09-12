package disaster

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateDisaster(t *testing.T) {
	tests := []struct {
		name         string
		req          CreateDisasterRequest
		repoErr      error
		wantErr      bool
		checkErr     error
		wantStartsAt *time.Time
	}{
		{
			name: "success",
			req: CreateDisasterRequest{
				Name:        "Assam Flood 2026",
				Description: "Severe flooding across districts",
				Type:        "flood",
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "normalizes timestamps to utc",
			req: CreateDisasterRequest{
				Name:     "Earthquake",
				Type:     "earthquake",
				StartsAt: utcPtr(time.Date(2026, 9, 6, 12, 30, 0, 0, time.FixedZone("IST", 5*3600+30*60))),
			},
			repoErr:      nil,
			wantErr:      false,
			wantStartsAt: utcPtr(time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)),
		},
		{
			name: "repository error",
			req: CreateDisasterRequest{
				Name: "Earthquake",
				Type: "earthquake",
			},
			repoErr:  assert.AnError,
			wantErr:  true,
			checkErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{createErr: tt.repoErr}
			svc := NewService(repo)

			resp, err := svc.CreateDisaster(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), resp.ID)
				assert.Equal(t, tt.req.Name, resp.Name)
				assert.Equal(t, tt.req.Type, resp.Type)
				assert.Equal(t, StatusActive, resp.Status)
				if tt.wantStartsAt != nil {
					if assert.NotNil(t, resp.StartsAt) {
						assert.Equal(t, tt.wantStartsAt, resp.StartsAt)
						assert.Equal(t, time.UTC, resp.StartsAt.Location())
					}
				}
			}
		})
	}
}

func utcPtr(t time.Time) *time.Time {
	return &t
}

func TestListDisasters(t *testing.T) {
	tests := []struct {
		name     string
		list     []Disaster
		repoErr  error
		wantErr  bool
		wantLen  int
		checkErr error
	}{
		{
			name: "success",
			list: []Disaster{
				{ID: 1, Name: "Assam Flood 2026", Type: "flood", Status: "active"},
				{ID: 2, Name: "Earthquake", Type: "earthquake", Status: "resolved"},
			},
			repoErr: nil,
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "empty",
			list:    []Disaster{},
			repoErr: nil,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:     "repository error",
			list:     nil,
			repoErr:  assert.AnError,
			wantErr:  true,
			checkErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{list: tt.list, listErr: tt.repoErr}
			svc := NewService(repo)

			resp, err := svc.ListDisasters(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp, tt.wantLen)
			}
		})
	}
}

func TestGetDisaster(t *testing.T) {
	tests := []struct {
		name     string
		get      *Disaster
		repoErr  error
		notFound bool
		wantErr  bool
		checkErr error
	}{
		{
			name:    "success",
			get:     &Disaster{ID: 1, Name: "Assam Flood 2026", Type: "flood", Status: "active"},
			wantErr: false,
		},
		{
			name:     "not found",
			notFound: true,
			wantErr:  true,
		},
		{
			name:     "repository error",
			repoErr:  assert.AnError,
			wantErr:  true,
			checkErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{get: tt.get, getErr: tt.repoErr, getNotFound: tt.notFound}
			svc := NewService(repo)

			resp, err := svc.GetDisaster(context.Background(), 1)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.get, resp)
			}
		})
	}
}
