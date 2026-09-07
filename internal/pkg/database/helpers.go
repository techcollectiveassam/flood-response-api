package database

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// NullableString converts an empty string to nil for nullable database columns.
func NullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// TextValue converts a nullable string column value to a plain string.
func TextValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// Timestamptz converts a *time.Time to a nullable timestamptz value.
func Timestamptz(value *time.Time) pgtype.Timestamptz {
	if value == nil || value.IsZero() {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

// FromTimestamptz converts a nullable timestamptz value to a *time.Time.
func FromTimestamptz(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}
