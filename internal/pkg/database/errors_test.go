package database

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
)

func TestTranslateError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "nil",
			err:        nil,
			wantStatus: 0,
			wantCode:   "",
		},
		{
			name:       "foreign key violation with known constraint",
			err:        &pgconn.PgError{Code: "23503", ConstraintName: "affected_areas_disaster_id_fkey"},
			wantStatus: 404,
			wantCode:   "disaster_not_found",
		},
		{
			name:       "foreign key violation with unknown constraint",
			err:        &pgconn.PgError{Code: "23503", ConstraintName: "unknown_fkey"},
			wantStatus: 409,
			wantCode:   "invalid_reference",
		},
		{
			name:       "not null violation",
			err:        &pgconn.PgError{Code: "23502"},
			wantStatus: 400,
			wantCode:   "missing_required_field",
		},
		{
			name:       "unique violation",
			err:        &pgconn.PgError{Code: "23505"},
			wantStatus: 409,
			wantCode:   "duplicate_record",
		},
		{
			name:       "string too long",
			err:        &pgconn.PgError{Code: "22001"},
			wantStatus: 400,
			wantCode:   "value_too_long",
		},
		{
			name:       "invalid text representation",
			err:        &pgconn.PgError{Code: "22P02"},
			wantStatus: 400,
			wantCode:   "invalid_value",
		},
		{
			name:       "unknown sqlstate passes through",
			err:        &pgconn.PgError{Code: "57014"},
			wantStatus: 500,
			wantCode:   "internal_error",
		},
		{
			name:       "non pg error passes through",
			err:        assert.AnError,
			wantStatus: 500,
			wantCode:   "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translated := TranslateError(tt.err)

			if tt.err == nil {
				assert.Nil(t, translated)
				return
			}

			assert.Equal(t, tt.wantStatus, apperror.HTTPStatus(translated))
			assert.Equal(t, tt.wantCode, apperror.Code(translated))
		})
	}
}
