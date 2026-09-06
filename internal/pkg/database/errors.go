package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

const foreignKeyViolation = "23503"

// constraintToError maps PostgreSQL constraint names to the API error they
// represent. Constraint violations carry one piece of context (the constraint
// name), so table-specific rules like "disaster no longer exists" live here.
var constraintToError = map[string]*apperror.Error{
	"affected_areas_disaster_id_fkey": apperror.NotFound("disaster_not_found", "disaster not found"),
}

func TranslateError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case foreignKeyViolation:
		if mapped := constraintToError[pgErr.ConstraintName]; mapped != nil {
			return mapped
		}
		return apperror.Conflict("invalid_reference", "the referenced record does not exist")
	case "23502":
		return apperror.BadRequest("missing_required_field", "a required field is missing")
	case "23505":
		return apperror.Conflict("duplicate_record", "a record with the same key already exists")
	case "22001":
		return apperror.BadRequest("value_too_long", "one or more field values are too long")
	case "22P02":
		return apperror.BadRequest("invalid_value", "one or more field values are invalid")
	default:
		return err
	}
}