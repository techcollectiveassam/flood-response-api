package affectedarea

import "time"

const (
	StatusReported      = "reported"
	StatusInvestigating = "investigating"
)

type AffectedArea struct {
	ID                 int64
	Name               string
	Description        string
	DisasterID         int32
	Location           string
	Geometry           string
	Latitude           *float64
	Longitude          *float64
	Severity           string
	VerificationStatus string
	ReportCount        int64
}

type AffectedAreaReport struct {
	ID              int64
	AreaID          int64
	Name            string
	Description     string
	LocationPayload []byte
	Severity        string
	ReporterName    string
	ReporterMobile  string
	CreatedAt       time.Time
}
