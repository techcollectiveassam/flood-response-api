package commandcenter

import "time"

const (
	TypeGovernment = "government"
	TypeNGO        = "ngo"
	TypeGroup      = "group"
	TypeOther      = "other"
)

type CommandCenter struct {
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
