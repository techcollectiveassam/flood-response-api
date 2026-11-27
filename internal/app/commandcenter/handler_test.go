package commandcenter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/command-centers", handler.CreateCommandCenter)
	router.GET("/disasters/:ref_id/command-center", handler.GetCommandCenterByDisaster)
	return router
}

func TestToCommandCenterResponse(t *testing.T) {
	disasterRefID := uuid.MustParse(testDisasterRefID)
	c := &CommandCenter{
		ID:            1,
		DisasterID:    2,
		DisasterRefID: disasterRefID,
		Name:          "Assam State Disaster Management",
		Type:          TypeGovernment,
		Description:   "State level coordination",
		ContactPerson: "Anjali Sharma",
		ContactMobile: "0000000000",
		ContactEmail:  "asdm@example.org",
	}

	resp := toCommandCenterResponse(c)

	assert.Equal(t, c.ID, resp.ID)
	assert.Equal(t, c.DisasterRefID.String(), resp.DisasterRefID)
	assert.Equal(t, c.Name, resp.Name)
	assert.Equal(t, c.Type, resp.Type)
	assert.Equal(t, c.Description, resp.Description)
	assert.Equal(t, c.ContactPerson, resp.ContactPerson)
	assert.Equal(t, c.ContactMobile, resp.ContactMobile)
	assert.Equal(t, c.ContactEmail, resp.ContactEmail)
}

func TestCreateCommandCenterHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		resolveErr error
		repoErr    error
		wantStatus int
		wantCode   string
		wantErr    bool
	}{
		{
			name: "success",
			body: CreateCommandCenterRequest{
				DisasterRefID: testDisasterRefID,
				Name:          "Assam State Disaster Management",
				Type:          TypeGovernment,
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
				"name": "Volunteer Group",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "invalid type",
			body: CreateCommandCenterRequest{
				DisasterRefID: testDisasterRefID,
				Name:          "Unknown",
				Type:          "corporation",
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "invalid disaster ref id",
			body: CreateCommandCenterRequest{
				DisasterRefID: "not-a-uuid",
				Name:          "Invalid Ref",
				Type:          TypeNGO,
			},
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "unknown disaster",
			body: CreateCommandCenterRequest{
				DisasterRefID: "99999999-9999-9999-9999-999999999999",
				Name:          "Orphan Response",
				Type:          TypeNGO,
			},
			resolveErr: ErrDisasterNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   "disaster_not_found",
			wantErr:    true,
		},
		{
			name: "service error",
			body: CreateCommandCenterRequest{
				DisasterRefID: "33333333-3333-3333-3333-333333333333",
				Name:          "Earthquake Response",
				Type:          TypeNGO,
			},
			repoErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{createErr: tt.repoErr, resolveErr: tt.resolveErr}
			svc := NewService(repo)
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

			req := httptest.NewRequest(http.MethodPost, "/command-centers", bytes.NewBuffer(body))
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
				assert.Equal(t, testDisasterRefID, resp.Data["disaster_ref_id"])
				assert.NotContains(t, resp.Data, "disaster_id")
				assert.Equal(t, "Assam State Disaster Management", resp.Data["name"])
				assert.Equal(t, "government", resp.Data["type"])
			}
		})
	}
}

func TestGetCommandCenterByDisasterHandler(t *testing.T) {
	successRefID := uuid.MustParse(testDisasterRefID)

	tests := []struct {
		name       string
		path       string
		repo       *mockRepository
		wantStatus int
		wantCode   string
		wantErr    bool
	}{
		{
			name: "success",
			path: "/disasters/" + testDisasterRefID + "/command-center",
			repo: &mockRepository{get: &CommandCenter{
				ID:            7,
				DisasterID:    2,
				Name:          "Volunteer Response Group",
				Type:          TypeGroup,
				Description:   "Flood relief coordination",
				ContactPerson: "Ravi Das",
				ContactMobile: "0000000000",
				ContactEmail:  "volunteer@example.org",
				CreatedAt:     mockCreatedAt,
				UpdatedAt:     mockCreatedAt,
			}},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "invalid ref id",
			path:       "/disasters/abc/command-center",
			repo:       &mockRepository{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_id",
			wantErr:    true,
		},
		{
			name:       "no command center for disaster",
			path:       "/disasters/42424242-4242-4242-4242-424242424242/command-center",
			repo:       &mockRepository{getErr: ErrCommandCenterNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "command_center_not_found",
			wantErr:    true,
		},
		{
			name:       "unknown disaster",
			path:       "/disasters/99999999-9999-9999-9999-999999999999/command-center",
			repo:       &mockRepository{resolveErr: ErrDisasterNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "disaster_not_found",
			wantErr:    true,
		},
		{
			name:       "service error",
			path:       "/disasters/33333333-3333-3333-3333-333333333333/command-center",
			repo:       &mockRepository{getErr: assert.AnError},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.repo)
			handler := NewHandler(svc)
			router := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
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
				assert.Equal(t, float64(7), resp.Data["id"])
				assert.Equal(t, successRefID.String(), resp.Data["disaster_ref_id"])
				assert.NotContains(t, resp.Data, "disaster_id")
				assert.Equal(t, "Volunteer Response Group", resp.Data["name"])
				assert.Equal(t, "group", resp.Data["type"])
				assert.Equal(t, "Flood relief coordination", resp.Data["description"])
				assert.Equal(t, "Ravi Das", resp.Data["contact_person"])
			}
		})
	}
}
