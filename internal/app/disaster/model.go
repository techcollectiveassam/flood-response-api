package disaster

import (
	"time"

	"github.com/google/uuid"
)

const StatusActive = "active"

const (
	TypeFlood      = "flood"
	TypeEarthquake = "earthquake"
	TypeLandslide  = "landslide"
)

type Disaster struct {
	ID          int32
	RefID       uuid.UUID
	Name        string
	Description string
	Type        string
	Status      string
	StartsAt    *time.Time
	EndsAt      *time.Time
}
