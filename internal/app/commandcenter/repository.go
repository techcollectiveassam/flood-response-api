package commandcenter

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

type Repository interface {
	Create(ctx context.Context, commandCenter *CommandCenter) error
}

type PostgresRepository struct {
	queries *sqlcgen.Queries
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, commandCenter *CommandCenter) error {
	result, err := r.queries.CreateCommandCenter(ctx, sqlcgen.CreateCommandCenterParams{
		Name:          commandCenter.Name,
		Type:          sqlcgen.CommandCenterType(commandCenter.Type),
		Description:   database.NullableString(commandCenter.Description),
		ContactPerson: database.NullableString(commandCenter.ContactPerson),
		ContactMobile: database.NullableString(commandCenter.ContactMobile),
		ContactEmail:  database.NullableString(commandCenter.ContactEmail),
	})
	if err != nil {
		return database.TranslateError(err)
	}
	commandCenter.ID = result.ID
	commandCenter.Name = result.Name
	commandCenter.Type = string(result.Type)
	commandCenter.Description = database.TextValue(result.Description)
	commandCenter.ContactPerson = database.TextValue(result.ContactPerson)
	commandCenter.ContactMobile = database.TextValue(result.ContactMobile)
	commandCenter.ContactEmail = database.TextValue(result.ContactEmail)
	commandCenter.CreatedAt = result.CreatedAt
	commandCenter.UpdatedAt = result.UpdatedAt
	return nil
}
