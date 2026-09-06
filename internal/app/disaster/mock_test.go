package disaster

import (
	"context"
	"time"
)

type mockRepository struct {
	createErr error
}

func (m *mockRepository) Create(ctx context.Context, d *Disaster) error {
	if m.createErr != nil {
		return m.createErr
	}
	d.ID = 1
	t := time.Now()
	d.Status = "active"
	d.StartsAt = &t
	return nil
}
