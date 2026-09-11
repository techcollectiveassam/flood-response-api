package app

import (
	"github.com/techcollectiveassam/flood-response-api/internal/app/affectedarea"
	"github.com/techcollectiveassam/flood-response-api/internal/app/disaster"
)

// Features holds the entry points for the application's feature modules.
type Features struct {
	AffectedArea *affectedarea.Module
	Disaster     *disaster.Module
}

func NewFeatures(application *Application) *Features {
	return &Features{
		AffectedArea: affectedarea.New(
			application.DB,
			application.Logger,
		),
		Disaster: disaster.New(
			application.DB,
		),
	}
}
