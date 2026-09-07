package disaster

import "time"

const StatusActive = "active"

const (
	TypeFlood      = "flood"
	TypeEarthquake = "earthquake"
	TypeLandslide  = "landslide"
)

type Disaster struct {
	ID          int32
	Name        string
	Description string
	Type        string
	Status      string
	StartsAt    *time.Time
	EndsAt      *time.Time
}
