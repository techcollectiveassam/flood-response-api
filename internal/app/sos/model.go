package sos

import "time"

const (
	StatusReported = "reported"
)

type SOS struct {
	ID             int64
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
