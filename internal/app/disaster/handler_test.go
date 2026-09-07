package disaster

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/disasters", handler.CreateDisaster)
	router.GET("/disasters", handler.ListDisasters)
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

func TestToDisasterResponseFormatsUTC(t *testing.T) {
	startsAt := time.Date(2026, 9, 6, 12, 30, 0, 0, time.FixedZone("IST", 5*3600+30*60))
	d := &Disaster{
		ID:       1,
		Name:     "Assam Flood 2026",
		Type:     "flood",
		Status:   "active",
		StartsAt: &startsAt,
	}

	resp := toDisasterResponse(d)

	if assert.NotNil(t, resp.StartsAt) {
		assert.Equal(t, "2026-09-06T07:00:00Z", *resp.StartsAt)
	}
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

func TestListDisastersHandler(t *testing.T) {
	tests := []struct {
		name       string
		list       []Disaster
		repoErr    error
		wantStatus int
		wantCode   string
		wantLen    int
		wantErr    bool
	}{
		{
			name: "success",
			list: []Disaster{
				{ID: 1, Name: "Assam Flood 2026", Type: "flood", Status: "active"},
				{ID: 2, Name: "Earthquake", Type: "earthquake", Status: "resolved"},
			},
			repoErr:    nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
			wantLen:    2,
		},
		{
			name:       "empty",
			list:       []Disaster{},
			repoErr:    nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
			wantLen:    0,
		},
		{
			name:       "service error",
			list:       nil,
			repoErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{list: tt.list, listErr: tt.repoErr}
			svc := NewService(repo, slog.Default())
			handler := NewHandler(svc)
			router := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/disasters", nil)
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
					Data []map[string]interface{} `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Data, tt.wantLen)
				if tt.wantLen > 0 {
					assert.Equal(t, "Assam Flood 2026", resp.Data[0]["name"])
					assert.Equal(t, "flood", resp.Data[0]["type"])
				}
			}
		})
	}
}
