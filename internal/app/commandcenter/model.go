package commandcenter

import (
	"time"

	"github.com/google/uuid"
)

const (
	TypeGovernment = "government"
	TypeNGO        = "ngo"
	TypeGroup      = "group"
	TypeOther      = "other"
)

type CommandCenter struct {
	DisasterID    int32
	DisasterRefID uuid.UUID
	ID            int32
	Name          string
	Type          string
	Description   string
	ContactPerson string
	ContactMobile string
	ContactEmail  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
