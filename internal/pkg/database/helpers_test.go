package database

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestNullableString(t *testing.T) {
	assert.Nil(t, NullableString(""))
	if v := NullableString("foo"); assert.NotNil(t, v) {
		assert.Equal(t, "foo", *v)
	}
}

func TestTextValue(t *testing.T) {
	assert.Equal(t, "", TextValue(nil))
	val := "foo"
	assert.Equal(t, "foo", TextValue(&val))
}

func TestTimestamptz(t *testing.T) {
	assert.False(t, Timestamptz(nil).Valid)
	assert.False(t, Timestamptz(&time.Time{}).Valid)

	ts := time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)
	v := Timestamptz(&ts)
	assert.True(t, v.Valid)
	assert.Equal(t, ts, v.Time)
}

func TestFromTimestamptz(t *testing.T) {
	assert.Nil(t, FromTimestamptz(pgtype.Timestamptz{}))

	ts := time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)
	got := FromTimestamptz(Timestamptz(&ts))
	if assert.NotNil(t, got) {
		assert.Equal(t, ts, *got)
	}
}
