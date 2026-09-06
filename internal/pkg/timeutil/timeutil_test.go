package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUTC(t *testing.T) {
	loc := time.FixedZone("IST", 5*3600+30*60)
	val := time.Date(2026, 9, 6, 12, 30, 0, 0, loc)
	got := UTC(&val)

	if assert.NotNil(t, got) {
		assert.Equal(t, time.UTC, got.Location())
		assert.Equal(t, time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC), *got)
	}

	assert.Nil(t, UTC(nil))
}

func TestFormatUTC(t *testing.T) {
	loc := time.FixedZone("IST", 5*3600+30*60)
	val := time.Date(2026, 9, 6, 12, 30, 0, 0, loc)
	assert.Equal(t, "2026-09-06T07:00:00Z", FormatUTC(val))

	utc := time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)
	assert.Equal(t, "2026-09-06T07:00:00Z", FormatUTC(utc))
}

func TestUTCZero(t *testing.T) {
	val := time.Time{}
	assert.Nil(t, UTC(&val))
}
