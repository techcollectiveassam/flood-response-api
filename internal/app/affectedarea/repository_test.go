package affectedarea

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

func TestPostgresRepositoryWithTxRejectsReentrancy(t *testing.T) {
	repo := &PostgresRepository{}

	err := repo.WithTx(context.Background(), func(Repository) error { return nil })

	assert.Error(t, err)
	assert.Equal(t, 500, apperror.HTTPStatus(err))
	assert.Equal(t, "nested_transaction", apperror.Code(err))
}
