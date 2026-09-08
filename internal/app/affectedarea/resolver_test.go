package affectedarea

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveCreatesNewAreaForAnySource(t *testing.T) {
	sources := []LocationPayload{
		{Source: SourceText, Description: "Village near the old bridge"},
		{Source: SourceGPS, Latitude: floatPointer(26.1445), Longitude: floatPointer(91.7362)},
		{Source: SourcePlaceSearch, Latitude: floatPointer(26.1204), Longitude: floatPointer(91.8006)},
		{Source: SourceMap, Geometry: []byte(`{"type":"Polygon","coordinates":[]}`)},
	}

	for _, loc := range sources {
		t.Run(loc.Source, func(t *testing.T) {
			resolver := NewAffectedAreaResolver()

			resolved, err := resolver.Resolve(context.Background(), 1, loc)

			assert.NoError(t, err)
			assert.Nil(t, resolved.Area)
			assert.Equal(t, DecisionNew, resolved.Decision)
			assert.Equal(t, 0.0, resolved.Confidence)
		})
	}
}

func TestResolveUnmatchableForText(t *testing.T) {
	resolver := NewAffectedAreaResolver()

	resolved, err := resolver.Resolve(context.Background(), 1, LocationPayload{
		Source:      SourceText,
		Description: "Village near the old bridge",
	})

	assert.NoError(t, err)
	assert.Equal(t, MatchReasonUnmatchable, resolved.MatchReason)
}
