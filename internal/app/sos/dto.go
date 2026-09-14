package sos

import (
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/timeutil"
)

type CreateSOSRequest struct {
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	ReporterMobile string   `json:"reporter_mobile"`
	Message        string   `json:"message"`
}

type SOSResponse struct {
	ID             int64    `json:"id"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	ReporterMobile string   `json:"reporter_mobile,omitempty"`
	Message        string   `json:"message,omitempty"`
	Status         string   `json:"status"`
	ReportCount    int32    `json:"report_count"`
	CreatedAt      string   `json:"created_at"`
}

func toSOSResponse(s *SOS) SOSResponse {
	return SOSResponse{
		ID:             s.ID,
		Latitude:       s.Latitude,
		Longitude:      s.Longitude,
		ReporterMobile: s.ReporterMobile,
		Message:        s.Message,
		Status:         s.Status,
		ReportCount:    s.ReportCount,
		CreatedAt:      timeutil.FormatUTC(s.CreatedAt),
	}
}
