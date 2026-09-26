package commandcenter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
				Name: "Volunteer Response Group",
				Type: TypeGroup,
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "repository error",
			req: CreateCommandCenterRequest{
				Name: "Earthquake Response",
				Type: TypeNGO,
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
