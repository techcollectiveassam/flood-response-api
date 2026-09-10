package affectedarea

import (
	"encoding/json"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/timeutil"
)

const (
	SourceGPS         = "gps"
	SourcePlaceSearch = "place_search"
	SourceMap         = "map"
	SourceText        = "text"
)

type LocationPayload struct {
	Source      string          `json:"source" binding:"required,oneof=gps place_search map text"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Latitude    *float64        `json:"latitude"`
	Longitude   *float64        `json:"longitude"`
	Geometry    json.RawMessage `json:"geometry"`
}

func (l LocationPayload) Validate() error {
	switch l.Source {
	case SourceGPS, SourcePlaceSearch:
		if l.Latitude == nil || l.Longitude == nil {
			return apperror.BadRequest("invalid_location", "latitude and longitude are required for "+l.Source+" source")
		}
	case SourceMap:
		if len(l.Geometry) == 0 {
			return apperror.BadRequest("invalid_location", "geometry is required for map source")
		}
		var value interface{}
		if err := json.Unmarshal(l.Geometry, &value); err != nil {
			return apperror.BadRequest("invalid_location", "geometry must be a valid GeoJSON object")
		}
	case SourceText:
		if l.Description == "" {
			return apperror.BadRequest("invalid_location", "description is required for text source")
		}
	}
	return nil
}

type Reporter struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

type CreateAffectedAreaRequest struct {
	DisasterID  int32            `json:"disaster_id" binding:"required"`
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	Location    *LocationPayload `json:"location" binding:"required"`
	Severity    string           `json:"severity" binding:"required,oneof=low medium high critical"`
	Reporter    *Reporter        `json:"reporter"`
}

type AffectedAreaResponse struct {
	ID                 int64    `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	DisasterID         int32    `json:"disaster_id"`
	Location           string   `json:"location,omitempty"`
	Geometry           string   `json:"geometry,omitempty"`
	Latitude           *float64 `json:"latitude,omitempty"`
	Longitude          *float64 `json:"longitude,omitempty"`
	Severity           string   `json:"severity"`
	VerificationStatus string   `json:"verification_status"`
	ReportCount        int64    `json:"report_count"`
}

type AffectedAreaReportResponse struct {
	ID        int64           `json:"id"`
	AreaID    int64           `json:"area_id"`
	Name      string          `json:"name"`
	Severity  string          `json:"severity"`
	Location  json.RawMessage `json:"location"`
	Reporter  *Reporter       `json:"reporter,omitempty"`
	CreatedAt string          `json:"created_at"`
}

type CreateAffectedAreaResponse struct {
	AffectedArea AffectedAreaResponse       `json:"affected_area"`
	Report       AffectedAreaReportResponse `json:"report"`
	IsNewArea    bool                       `json:"is_new_area"`
	Confidence   float64                    `json:"confidence"`
	MatchReason  string                     `json:"match_reason"`
}

func toAffectedAreaResponse(area *AffectedArea) AffectedAreaResponse {
	return AffectedAreaResponse{
		ID:                 area.ID,
		Name:               area.Name,
		Description:        area.Description,
		DisasterID:         area.DisasterID,
		Location:           area.Location,
		Geometry:           area.Geometry,
		Latitude:           area.Latitude,
		Longitude:          area.Longitude,
		Severity:           area.Severity,
		VerificationStatus: area.VerificationStatus,
		ReportCount:        area.ReportCount,
	}
}

func toAffectedAreaReportResponse(report *AffectedAreaReport) AffectedAreaReportResponse {
	resp := AffectedAreaReportResponse{
		ID:        report.ID,
		AreaID:    report.AreaID,
		Name:      report.Name,
		Severity:  report.Severity,
		Location:  report.LocationPayload,
		CreatedAt: timeutil.FormatUTC(report.CreatedAt),
	}
	if len(resp.Location) == 0 {
		resp.Location = json.RawMessage(`{}`)
	}
	if report.ReporterName != "" || report.ReporterMobile != "" {
		resp.Reporter = &Reporter{Name: report.ReporterName, Mobile: report.ReporterMobile}
	}
	return resp
}

func reporterName(reporter *Reporter) string {
	if reporter == nil {
		return ""
	}
	return reporter.Name
}

func reporterMobile(reporter *Reporter) string {
	if reporter == nil {
		return ""
	}
	return reporter.Mobile
}

func locationLabel(loc LocationPayload) string {
	switch loc.Source {
	case SourceText:
		return loc.Description
	default:
		return loc.Name
	}
}

func locationGeometry(loc LocationPayload) string {
	if loc.Source == SourceMap {
		return string(loc.Geometry)
	}
	return ""
}

func locationPayloadJSON(loc LocationPayload) []byte {
	payload, err := json.Marshal(loc)
	if err != nil {
		return []byte(`{}`)
	}
	return payload
}
