package affectedarea

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAffectedArea(t *testing.T) {
	logger := slog.Default()

	tests := []struct {
		name     string
		req      CreateAffectedAreaRequest
		repoErr  error
		wantErr  bool
		checkErr error
	}{
		{
			name: "success",
			req: CreateAffectedAreaRequest{
				Name:        "Flood Zone A",
				DisasterID:  "disaster-123",
				Severity:    "high",
				Description: "Severe flooding",
				Location:    "26.1445,91.7362",
				Source:      "satellite",
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "repository error",
			req: CreateAffectedAreaRequest{
				Name:       "Flood Zone B",
				DisasterID: "disaster-456",
				Severity:   "medium",
			},
			repoErr:  ErrNotImplemented,
			wantErr:  true,
			checkErr: ErrNotImplemented,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{createErr: tt.repoErr}
			svc := NewService(repo, logger)

			err := svc.CreateAffectedArea(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.checkErr != nil {
					assert.ErrorIs(t, err, tt.checkErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
