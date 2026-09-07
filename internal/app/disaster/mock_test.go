package disaster

import (
	"context"
	"time"
)

var mockStartsAt = time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)

type mockRepository struct {
	createErr error
	listErr   error
	list      []Disaster
}

func (m *mockRepository) Create(ctx context.Context, d *Disaster) error {
	if m.createErr != nil {
		return m.createErr
	}
	d.ID = 1
	d.Status = "active"
	if d.StartsAt == nil {
		t := mockStartsAt
		d.StartsAt = &t
	}
	return nil
}

func (m *mockRepository) List(ctx context.Context) ([]Disaster, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.list, nil
}
