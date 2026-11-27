package commandcenter

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

const testDisasterRefID = "11111111-1111-1111-1111-111111111111"

func TestCreateCommandCenter(t *testing.T) {
	tests := []struct {
		name     string
		req      CreateCommandCenterRequest
		resolve  error
		repoErr  error
		wantErr  bool
		checkErr error
	}{
		{
			name: "success",
			req: CreateCommandCenterRequest{
				DisasterRefID: testDisasterRefID,
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
				DisasterRefID: "22222222-2222-2222-2222-222222222222",
				Name:          "Volunteer Response Group",
				Type:          TypeGroup,
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "invalid disaster ref id",
			req: CreateCommandCenterRequest{
				DisasterRefID: "not-a-uuid",
				Name:          "Invalid Ref",
				Type:          TypeNGO,
			},
			wantErr:  true,
			checkErr: ErrInvalidDisasterRef,
		},
		{
			name: "unknown disaster",
			req: CreateCommandCenterRequest{
				DisasterRefID: "99999999-9999-9999-9999-999999999999",
				Name:          "Orphan Response",
				Type:          TypeNGO,
			},
			resolve:  ErrDisasterNotFound,
			wantErr:  true,
			checkErr: ErrDisasterNotFound,
		},
		{
			name: "repository error",
			req: CreateCommandCenterRequest{
				DisasterRefID: "33333333-3333-3333-3333-333333333333",
				Name:          "Earthquake Response",
				Type:          TypeNGO,
			},
			repoErr:  assert.AnError,
			wantErr:  true,
			checkErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{createErr: tt.repoErr, resolveErr: tt.resolve}
			svc := NewService(repo)

			resp, err := svc.CreateCommandCenter(context.Background(), tt.req)

			if tt.checkErr == ErrInvalidDisasterRef {
				assert.Equal(t, 0, repo.resolveCalls)
			} else {
				assert.Equal(t, 1, repo.resolveCalls)
				assert.Equal(t, tt.req.DisasterRefID, repo.resolvedRefID.String())
			}

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), resp.ID)
				assert.Equal(t, int32(7), resp.DisasterID)
				assert.Equal(t, tt.req.DisasterRefID, resp.DisasterRefID.String())
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
	refID := uuid.MustParse(testDisasterRefID)
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
		refID        uuid.UUID
		repo         *mockRepository
		wantErr      bool
		checkErr     error
		wantResolve  int
		wantGetCalls int
	}{
		{
			name:         "success",
			refID:        refID,
			repo:         &mockRepository{get: found},
			wantErr:      false,
			wantResolve:  1,
			wantGetCalls: 1,
		},
		{
			name:         "disaster exists but has no command center",
			refID:        uuid.MustParse("42424242-4242-4242-4242-424242424242"),
			repo:         &mockRepository{getErr: ErrCommandCenterNotFound},
			wantErr:      true,
			checkErr:     ErrCommandCenterNotFound,
			wantResolve:  1,
			wantGetCalls: 1,
		},
		{
			name:         "unknown disaster",
			refID:        uuid.MustParse("99999999-9999-9999-9999-999999999999"),
			repo:         &mockRepository{resolveErr: ErrDisasterNotFound},
			wantErr:      true,
			checkErr:     ErrDisasterNotFound,
			wantResolve:  1,
			wantGetCalls: 0,
		},
		{
			name:         "resolution fails",
			refID:        uuid.MustParse("55555555-5555-5555-5555-555555555555"),
			repo:         &mockRepository{resolveErr: assert.AnError},
			wantErr:      true,
			checkErr:     assert.AnError,
			wantResolve:  1,
			wantGetCalls: 0,
		},
		{
			name:         "repository error",
			refID:        uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			repo:         &mockRepository{getErr: assert.AnError},
			wantErr:      true,
			checkErr:     assert.AnError,
			wantResolve:  1,
			wantGetCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.repo)

			resp, err := svc.GetCommandCenterByDisaster(context.Background(), tt.refID)

			assert.Equal(t, tt.refID, tt.repo.resolvedRefID)
			assert.Equal(t, tt.wantResolve, tt.repo.resolveCalls)
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
				assert.Equal(t, tt.refID.String(), resp.DisasterRefID.String())
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
