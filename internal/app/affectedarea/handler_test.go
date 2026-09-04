package affectedarea

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/affected-areas", handler.CreateAffectedArea)
	return router
}

func TestCreateAffectedAreaHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		repoErr    error
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone A",
				DisasterID: "disaster-123",
				Severity:   "high",
			},
			repoErr:    nil,
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:       "invalid json",
			body:       "invalid",
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "missing required fields",
			body: map[string]string{
				"name": "Flood Zone A",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "service error",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone B",
				DisasterID: "disaster-456",
				Severity:   "medium",
			},
			repoErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "not implemented",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone C",
				DisasterID: "disaster-789",
				Severity:   "low",
			},
			repoErr:    ErrNotImplemented,
			wantStatus: http.StatusNotImplemented,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{createErr: tt.repoErr}
			svc := NewService(repo, slog.Default())
			handler := NewHandler(svc)
			router := setupRouter(handler)

			var body []byte
			var err error
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/affected-areas", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantErr {
				var resp map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp["error"])
			}
		})
	}
}
