package sos

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusReported = "reported"
)

type SOS struct {
	ID             int64
	TrackingID     uuid.UUID
	DisasterID     int32
	Latitude       *float64
	Longitude      *float64
	ReporterMobile string
	Message        string
	Status         string
	ReportCount    int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
