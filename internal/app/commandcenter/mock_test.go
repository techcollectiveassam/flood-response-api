package commandcenter

import (
	"context"
	"time"
)

var mockCreatedAt = time.Date(2026, 9, 25, 7, 0, 0, 0, time.UTC)

type mockRepository struct {
	createErr error
	create    *CommandCenter
}

func (m *mockRepository) Create(ctx context.Context, c *CommandCenter) error {
	if m.createErr != nil {
		return m.createErr
	}
	c.ID = 1
	c.CreatedAt = mockCreatedAt
	c.UpdatedAt = mockCreatedAt
	return nil
}
