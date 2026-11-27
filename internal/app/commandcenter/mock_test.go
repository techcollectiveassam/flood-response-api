package commandcenter

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var mockCreatedAt = time.Date(2026, 9, 25, 7, 0, 0, 0, time.UTC)

type mockRepository struct {
	createErr   error
	create      *CommandCenter
	getErr      error
	get         *CommandCenter
	getCalls    int
	gotDisaster int32
	resolveErr  error

	resolveDisasterID int32
	resolveCalls      int
	resolvedRefID     uuid.UUID
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

func (m *mockRepository) GetByDisasterID(ctx context.Context, disasterID int32) (*CommandCenter, error) {
	m.getCalls++
	m.gotDisaster = disasterID
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.get, nil
}

func (m *mockRepository) GetDisasterIDByRefID(ctx context.Context, refID uuid.UUID) (int32, error) {
	m.resolveCalls++
	m.resolvedRefID = refID
	if m.resolveErr != nil {
		return 0, m.resolveErr
	}
	if m.resolveDisasterID == 0 {
		return 7, nil
	}
	return m.resolveDisasterID, nil
}
