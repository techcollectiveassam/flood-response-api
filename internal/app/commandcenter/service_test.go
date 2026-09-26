package commandcenter

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

func TestCreateCommandCenter(t *testing.T) {
	tests := []struct {
		name     string
		req      CreateCommandCenterRequest
		repoErr  error
		wantErr  bool
		checkErr error
	}{
		{
			name: "success",
			req: CreateCommandCenterRequest{
				DisasterID:    1,
				Name:          "Assam State Disaster Management",
				Type:          TypeGovernment,
				Description:   "State level coordination",
				ContactPerson: "Anjali Sharma",
				ContactMobile: "0000000000",
				ContactEmail:  "asdm@example.org",
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "success with only required fields",
			req: CreateCommandCenterRequest{
				DisasterID: 2,
				Name:       "Volunteer Response Group",
				Type:       TypeGroup,
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "repository error",
			req: CreateCommandCenterRequest{
				DisasterID: 3,
				Name:       "Earthquake Response",
				Type:       TypeNGO,
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

			resp, err := svc.CreateCommandCenter(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), resp.ID)
				assert.Equal(t, tt.req.DisasterID, resp.DisasterID)
				assert.Equal(t, tt.req.Name, resp.Name)
				assert.Equal(t, tt.req.Type, resp.Type)
				assert.Equal(t, tt.req.Description, resp.Description)
				assert.Equal(t, tt.req.ContactPerson, resp.ContactPerson)
				assert.Equal(t, tt.req.ContactMobile, resp.ContactMobile)
				assert.Equal(t, tt.req.ContactEmail, resp.ContactEmail)
				assert.Equal(t, mockCreatedAt, resp.CreatedAt)
				assert.Equal(t, mockCreatedAt, resp.UpdatedAt)
			}
		})
	}
}

func TestGetCommandCenterByDisaster(t *testing.T) {
	found := &CommandCenter{
		ID:            7,
		DisasterID:    2,
		Name:          "Volunteer Response Group",
		Type:          TypeGroup,
		ContactPerson: "Ravi Das",
		CreatedAt:     mockCreatedAt,
		UpdatedAt:     mockCreatedAt,
	}

	tests := []struct {
		name         string
		disasterID   int32
		repo         *mockRepository
		wantErr      bool
		checkErr     error
		wantExists   int
		wantGetCalls int
	}{
		{
			name:         "success",
			disasterID:   2,
			repo:         &mockRepository{get: found, disasterExists: true},
			wantErr:      false,
			wantExists:   1,
			wantGetCalls: 1,
		},
		{
			name:         "disaster exists but has no command center",
			disasterID:   4242,
			repo:         &mockRepository{getErr: ErrCommandCenterNotFound, disasterExists: true},
			wantErr:      true,
			checkErr:     ErrCommandCenterNotFound,
			wantExists:   1,
			wantGetCalls: 1,
		},
		{
			name:         "unknown disaster",
			disasterID:   999999,
			repo:         &mockRepository{disasterExists: false},
			wantErr:      true,
			checkErr:     ErrDisasterNotFound,
			wantExists:   1,
			wantGetCalls: 0,
		},
		{
			name:         "existence check fails",
			disasterID:   5555,
			repo:         &mockRepository{disasterExistsErr: assert.AnError},
			wantErr:      true,
			checkErr:     assert.AnError,
			wantExists:   1,
			wantGetCalls: 0,
		},
		{
			name:         "repository error",
			disasterID:   3,
			repo:         &mockRepository{getErr: assert.AnError, disasterExists: true},
			wantErr:      true,
			checkErr:     assert.AnError,
			wantExists:   1,
			wantGetCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.repo)

			resp, err := svc.GetCommandCenterByDisaster(context.Background(), tt.disasterID)

			assert.Equal(t, tt.disasterID, tt.repo.gotDisaster)
			assert.Equal(t, tt.wantExists, tt.repo.existsCalls)
			assert.Equal(t, tt.wantGetCalls, tt.repo.getCalls)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, found.ID, resp.ID)
				assert.Equal(t, found.DisasterID, resp.DisasterID)
				assert.Equal(t, found.Name, resp.Name)
				assert.Equal(t, found.Type, resp.Type)
			}
		})
	}
}

func TestErrCommandCenterNotFoundIs404(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, apperror.HTTPStatus(ErrCommandCenterNotFound))
	assert.Equal(t, "command_center_not_found", apperror.Code(ErrCommandCenterNotFound))
}

func TestErrDisasterNotFoundIs404(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, apperror.HTTPStatus(ErrDisasterNotFound))
	assert.Equal(t, "disaster_not_found", apperror.Code(ErrDisasterNotFound))
}
