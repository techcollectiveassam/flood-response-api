package affectedarea

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusReported      = "reported"
	StatusInvestigating = "investigating"
)

type AffectedArea struct {
	ID                 int64
	RefID              uuid.UUID
	Description        string
	DisasterID         int32
	DisasterRefID      uuid.UUID
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
	RefID           uuid.UUID
	AreaID          int64
	Description     string
	LocationPayload []byte
	Severity        string
	ReporterName    string
	ReporterMobile  string
	CreatedAt       time.Time
}
