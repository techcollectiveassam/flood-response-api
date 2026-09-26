package commandcenter

import (
	"context"
	"time"
)

var mockCreatedAt = time.Date(2026, 9, 25, 7, 0, 0, 0, time.UTC)

type mockRepository struct {
	createErr         error
	create            *CommandCenter
	getErr            error
	get               *CommandCenter
	getCalls          int
	gotDisaster       int32
	disasterExists    bool
	disasterExistsErr error
	existsCalls       int
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

func (m *mockRepository) DisasterExists(ctx context.Context, disasterID int32) (bool, error) {
	m.existsCalls++
	m.gotDisaster = disasterID
	return m.disasterExists, m.disasterExistsErr
}
