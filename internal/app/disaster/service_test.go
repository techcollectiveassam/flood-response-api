package disaster

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDisaster(t *testing.T) {
	logger := slog.Default()

	tests := []struct {
		name     string
		req      CreateDisasterRequest
		repoErr  error
		wantErr  bool
		checkErr error
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
			svc := NewService(repo, logger)

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
			}
		})
	}
}
