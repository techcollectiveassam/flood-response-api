package affectedarea

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/pagination"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/validation"
)

func setupRouter(handler *Handler, params pagination.Params) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/affected-areas", handler.CreateAffectedArea)
	router.GET("/affected-areas", pagination.Middleware(params), handler.ListAffectedAreas)
	return router
}

func testParams() pagination.Params {
	return pagination.Params{DefaultPage: 1, DefaultLimit: 20, MaxLimit: 100}
}

func registerJSONTagNames(t *testing.T) {
	t.Helper()
	v, ok := binding.Validator.Engine().(*validator.Validate)
	require.True(t, ok)
	validation.RegisterJSONTagNames(v)
}

func TestCreateAffectedAreaValidationDetails(t *testing.T) {
	registerJSONTagNames(t)
	svc := newTestService(&mockRepository{})
	handler := NewHandler(svc)
	router := setupRouter(handler, pagination.Params{})

	req := httptest.NewRequest(http.MethodPost, "/affected-areas", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp struct {
		Error struct {
			Code    string            `json:"code"`
			Message string            `json:"message"`
			Details []apperror.Detail `json:"details"`
		} `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_request_body", resp.Error.Code)
	assert.Equal(t, "validation error", resp.Error.Message)
	assert.Equal(t, []apperror.Detail{
		{Field: "disaster_id", Message: "disaster_id is required"},
		{Field: "name", Message: "name is required"},
		{Field: "location", Message: "location is required"},
		{Field: "severity", Message: "severity is required"},
	}, resp.Error.Details)
}

func TestToAffectedAreaResponse(t *testing.T) {
	latitude := 26.1445
	longitude := 91.7362
	area := &AffectedArea{
		ID:                 1,
		Name:               "Flood Zone A",
		Description:        "Severe flooding",
		DisasterID:         4,
		Location:           "Beltola, Guwahati",
		Latitude:           &latitude,
		Longitude:          &longitude,
		Severity:           "high",
		VerificationStatus: StatusReported,
		ReportCount:        3,
	}

	resp := toAffectedAreaResponse(area)

	assert.Equal(t, area.ID, resp.ID)
	assert.Equal(t, area.Name, resp.Name)
	assert.Equal(t, area.Description, resp.Description)
	assert.Equal(t, area.DisasterID, resp.DisasterID)
	assert.Equal(t, area.Location, resp.Location)
	assert.Equal(t, area.Latitude, resp.Latitude)
	assert.Equal(t, area.Longitude, resp.Longitude)
	assert.Equal(t, area.Severity, resp.Severity)
	assert.Equal(t, area.VerificationStatus, resp.VerificationStatus)
	assert.Equal(t, area.ReportCount, resp.ReportCount)
}

func TestCreateAffectedAreaHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		repo       *mockRepository
		wantStatus int
		wantCode   string
		wantErr    bool
	}{
		{
			name:       "success",
			body:       textRequest(),
			repo:       &mockRepository{},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:       "invalid json",
			body:       "invalid",
			repo:       &mockRepository{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "missing required fields",
			body: map[string]string{
				"name": "Flood Zone A",
			},
			repo:       &mockRepository{},
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
			repo:       &mockRepository{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request_body",
			wantErr:    true,
		},
		{
			name: "invalid gps location",
			body: CreateAffectedAreaRequest{
				Name:       "Flood Zone C",
				DisasterID: 1,
				Severity:   "medium",
				Location: &LocationPayload{
					Source:    SourceGPS,
					Latitude:  floatPointer(26.1445),
					Longitude: nil,
				},
			},
			repo:       &mockRepository{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_location",
			wantErr:    true,
		},
		{
			name: "disaster not found",
			body: textRequest(),
			repo: &mockRepository{
				createAreaErr: ErrDisasterNotFound,
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "disaster_not_found",
			wantErr:    true,
		},
		{
			name: "service error",
			body: textRequest(),
			repo: &mockRepository{
				createAreaErr: assert.AnError,
			},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			handler := NewHandler(svc)
			router := setupRouter(handler, pagination.Params{})

			var body []byte
			var err error
			switch v := tt.body.(type) {
			case CreateAffectedAreaRequest:
				body, err = json.Marshal(v)
				assert.NoError(t, err)
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
					Data struct {
						AffectedArea struct {
							ID   int64  `json:"id"`
							Name string `json:"name"`
						} `json:"affected_area"`
						Report struct {
							ID     int64  `json:"id"`
							AreaID int64  `json:"area_id"`
							Name   string `json:"name"`
						} `json:"report"`
						IsNewArea   bool    `json:"is_new_area"`
						Confidence  float64 `json:"confidence"`
						MatchReason string  `json:"match_reason"`
					} `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), resp.Data.AffectedArea.ID)
				assert.Equal(t, "Possible flooding", resp.Data.AffectedArea.Name)
				assert.Equal(t, int64(1), resp.Data.Report.ID)
				assert.Equal(t, int64(1), resp.Data.Report.AreaID)
				assert.True(t, resp.Data.IsNewArea)
				assert.Equal(t, MatchReasonUnmatchable, resp.Data.MatchReason)
			}
		})
	}
}

func TestListAffectedAreasHandler(t *testing.T) {
	area := &AffectedArea{
		ID:                 1,
		Name:               "Flood Zone A",
		Description:        "Severe flooding",
		DisasterID:         4,
		Location:           "Beltola, Guwahati",
		Severity:           "high",
		VerificationStatus: StatusReported,
		ReportCount:        3,
	}

	tests := []struct {
		name          string
		query         string
		repo          *mockRepository
		wantStatus    int
		wantErr       bool
		wantPage      int
		wantLimit     int
		wantTotal     int64
		wantPages     int
		wantItemCount int
	}{
		{
			name:          "default pagination",
			query:         "",
			repo:          &mockRepository{listAreas: []*AffectedArea{area}, listTotal: 57},
			wantStatus:    http.StatusOK,
			wantPage:      testParams().DefaultPage,
			wantLimit:     testParams().DefaultLimit,
			wantTotal:     57,
			wantPages:     3,
			wantItemCount: 1,
		},
		{
			name:          "explicit pagination",
			query:         "?page=2&limit=10",
			repo:          &mockRepository{listAreas: []*AffectedArea{area}, listTotal: 25},
			wantStatus:    http.StatusOK,
			wantPage:      2,
			wantLimit:     10,
			wantTotal:     25,
			wantPages:     3,
			wantItemCount: 1,
		},
		{
			name:          "invalid page",
			query:         "?page=0",
			repo:          &mockRepository{},
			wantStatus:    http.StatusBadRequest,
			wantErr:       true,
			wantItemCount: 0,
		},
		{
			name:          "limit exceeds max",
			query:         "?limit=101",
			repo:          &mockRepository{},
			wantStatus:    http.StatusBadRequest,
			wantErr:       true,
			wantItemCount: 0,
		},
		{
			name:          "service error",
			query:         "",
			repo:          &mockRepository{listErr: assert.AnError},
			wantStatus:    http.StatusInternalServerError,
			wantErr:       true,
			wantItemCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			handler := NewHandler(svc)
			router := setupRouter(handler, testParams())

			req := httptest.NewRequest(http.MethodGet, "/affected-areas"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantErr {
				var resp struct {
					Error struct {
						Code string `json:"code"`
					} `json:"error"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Error.Code)
				return
			}

			var resp struct {
				Data []struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"data"`
				Pagination struct {
					Page       int   `json:"page"`
					Limit      int   `json:"limit"`
					Total      int64 `json:"total"`
					TotalPages int   `json:"total_pages"`
				} `json:"pagination"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Len(t, resp.Data, tt.wantItemCount)
			if tt.wantItemCount > 0 {
				assert.Equal(t, int64(1), resp.Data[0].ID)
				assert.Equal(t, "Flood Zone A", resp.Data[0].Name)
			}
			assert.Equal(t, tt.wantPage, resp.Pagination.Page)
			assert.Equal(t, tt.wantLimit, resp.Pagination.Limit)
			assert.Equal(t, tt.wantTotal, resp.Pagination.Total)
			assert.Equal(t, tt.wantPages, resp.Pagination.TotalPages)
		})
	}
}
