package sos

import (
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/timeutil"
)

type CreateSOSRequest struct {
	DisasterID     int32     `json:"disaster_id" binding:"required,gt=0"`
	ReporterMobile string    `json:"reporter_mobile"`
	Message        string    `json:"message"`
	Location       *Location `json:"location"`
}

type Location struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type SOSResponse struct {
	ID             int64     `json:"id"`
	DisasterID     int32     `json:"disaster_id"`
	Location       *Location `json:"location,omitempty"`
	ReporterMobile string    `json:"reporter_mobile,omitempty"`
	Message        string    `json:"message,omitempty"`
	Status         string    `json:"status"`
	ReportCount    int32     `json:"report_count"`
	CreatedAt      string    `json:"created_at"`
}

func toSOSResponse(s *SOS) SOSResponse {
	resp := SOSResponse{
		ID:             s.ID,
		DisasterID:     s.DisasterID,
		ReporterMobile: s.ReporterMobile,
		Message:        s.Message,
		Status:         s.Status,
		ReportCount:    s.ReportCount,
		CreatedAt:      timeutil.FormatUTC(s.CreatedAt),
	}
	if s.Latitude != nil && s.Longitude != nil {
		resp.Location = &Location{Latitude: s.Latitude, Longitude: s.Longitude}
	}
	return resp
}
