package commandcenter

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

var (
	ErrDisasterNotFound      = apperror.NotFound("disaster_not_found", "disaster not found")
	ErrCommandCenterNotFound = apperror.NotFound("command_center_not_found", "command center not found")
)

type Repository interface {
	Create(ctx context.Context, commandCenter *CommandCenter) error
	GetByDisasterID(ctx context.Context, disasterID int32) (*CommandCenter, error)
	DisasterExists(ctx context.Context, disasterID int32) (bool, error)
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
		DisasterID:    commandCenter.DisasterID,
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
	commandCenter.DisasterID = result.DisasterID
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

func (r *PostgresRepository) GetByDisasterID(ctx context.Context, disasterID int32) (*CommandCenter, error) {
	row, err := r.queries.GetCommandCenterByDisaster(ctx, disasterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCommandCenterNotFound
		}
		return nil, database.TranslateError(err)
	}

	commandCenter := commandCenterFromRow(
		row.ID,
		row.DisasterID,
		row.Name,
		row.Type,
		row.Description,
		row.ContactPerson,
		row.ContactMobile,
		row.ContactEmail,
		row.CreatedAt,
		row.UpdatedAt,
	)
	return &commandCenter, nil
}

// DisasterExists reports whether the parent disaster of a nested command center
// lookup exists, so the service can tell an unknown disaster from a disaster
// that simply has no command center yet.
func (r *PostgresRepository) DisasterExists(ctx context.Context, disasterID int32) (bool, error) {
	exists, err := r.queries.DisasterExists(ctx, disasterID)
	if err != nil {
		return false, database.TranslateError(err)
	}
	return exists, nil
}

func commandCenterFromRow(
	id int32,
	disasterID int32,
	name string,
	ctype sqlcgen.CommandCenterType,
	description *string,
	contactPerson *string,
	contactMobile *string,
	contactEmail *string,
	createdAt time.Time,
	updatedAt time.Time,
) CommandCenter {
	return CommandCenter{
		ID:            id,
		DisasterID:    disasterID,
		Name:          name,
		Type:          string(ctype),
		Description:   database.TextValue(description),
		ContactPerson: database.TextValue(contactPerson),
		ContactMobile: database.TextValue(contactMobile),
		ContactEmail:  database.TextValue(contactEmail),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
