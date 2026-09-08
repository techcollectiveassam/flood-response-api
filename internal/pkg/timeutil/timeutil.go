package timeutil

import "time"

const UTCFormat = time.RFC3339

// UTC returns a UTC-normalized copy of t, or nil if t is nil or zero.
func UTC(t *time.Time) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}
	utc := t.UTC()
	return &utc
}

// FormatUTC returns t normalized to UTC and formatted as RFC3339 (always ends in Z).
func FormatUTC(t time.Time) string {
	return t.UTC().Format(UTCFormat)
}
