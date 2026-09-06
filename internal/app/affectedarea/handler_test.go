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
		wantCode   string
		wantErr    bool
	}{
		{
			name: "success",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone A",
				DisasterID: 1,
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
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "missing required fields",
			body: map[string]string{
				"name": "Flood Zone A",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "invalid disaster id",
			body: map[string]interface{}{
				"name":        "Flood Zone B",
				"disaster_id": "abc",
				"severity":    "medium",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "invalid severity",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone B",
				DisasterID: 1,
				Severity:   "urgent",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "disaster not found",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone B",
				DisasterID: 999,
				Severity:   "medium",
			},
			repoErr:    ErrDisasterNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   "disaster_not_found",
			wantErr:    true,
		},
		{
			name: "service error",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone B",
				DisasterID: 1,
				Severity:   "medium",
			},
			repoErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
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
				var resp struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Error.Code)
				assert.NotEmpty(t, resp.Error.Message)
				if tt.wantCode != "" {
					assert.Equal(t, tt.wantCode, resp.Error.Code)
				}
			} else {
				var resp struct {
					Data map[string]interface{} `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "1", resp.Data["id"])
				assert.Equal(t, "Flood Zone A", resp.Data["name"])
				assert.Equal(t, float64(1), resp.Data["disaster_id"])
				assert.Equal(t, "high", resp.Data["severity"])
			}
		})
	}
}
