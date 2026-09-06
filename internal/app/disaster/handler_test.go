package disaster

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
	router.POST("/disasters", handler.CreateDisaster)
	return router
}

func TestToDisasterResponse(t *testing.T) {
	d := &Disaster{
		ID:          1,
		Name:        "Assam Flood 2026",
		Description: "Severe flooding",
		Type:        "flood",
		Status:      "active",
	}

	resp := toDisasterResponse(d)

	assert.Equal(t, d.ID, resp.ID)
	assert.Equal(t, d.Name, resp.Name)
	assert.Equal(t, d.Description, resp.Description)
	assert.Equal(t, d.Type, resp.Type)
	assert.Equal(t, d.Status, resp.Status)
	assert.Nil(t, resp.StartsAt)
	assert.Nil(t, resp.EndsAt)
}

func TestCreateDisasterHandler(t *testing.T) {
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
			body: CreateDisasterRequest{
				Name: "Assam Flood 2026",
				Type: "flood",
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
				"name": "Flood",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "invalid type",
			body: CreateDisasterRequest{
				Name: "Cyclone",
				Type: "cyclone",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "service error",
			body: CreateDisasterRequest{
				Name: "Earthquake",
				Type: "earthquake",
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

			req := httptest.NewRequest(http.MethodPost, "/disasters", bytes.NewBuffer(body))
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
				assert.Equal(t, float64(1), resp.Data["id"])
				assert.Equal(t, "Assam Flood 2026", resp.Data["name"])
				assert.Equal(t, "flood", resp.Data["type"])
				assert.Equal(t, "active", resp.Data["status"])
			}
		})
	}
}
