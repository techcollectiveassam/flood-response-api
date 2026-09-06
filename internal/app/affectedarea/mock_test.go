package affectedarea

import "context"

type mockRepository struct {
	createErr error
}

func (m *mockRepository) Create(ctx context.Context, area *AffectedArea) error {
	if m.createErr != nil {
		return m.createErr
	}
	area.ID = "1"
	return nil
}
