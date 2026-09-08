package affectedarea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationPayloadValidate(t *testing.T) {
	tests := []struct {
		name    string
		payload LocationPayload
		wantErr bool
	}{
		{
			name: "gps with coordinates",
			payload: LocationPayload{
				Source:    SourceGPS,
				Latitude:  floatPointer(26.1445),
				Longitude: floatPointer(91.7362),
			},
			wantErr: false,
		},
		{
			name: "gps missing longitude",
			payload: LocationPayload{
				Source:    SourceGPS,
				Latitude:  floatPointer(26.1445),
				Longitude: nil,
			},
			wantErr: true,
		},
		{
			name: "place search with coordinates",
			payload: LocationPayload{
				Source:    SourcePlaceSearch,
				Name:      "Beltola, Guwahati",
				Latitude:  floatPointer(26.1204),
				Longitude: floatPointer(91.8006),
			},
			wantErr: false,
		},
		{
			name: "map with geometry",
			payload: LocationPayload{
				Source:   SourceMap,
				Geometry: []byte(`{"type":"Polygon","coordinates":[]}`),
			},
			wantErr: false,
		},
		{
			name: "map with invalid geometry",
			payload: LocationPayload{
				Source:   SourceMap,
				Geometry: []byte(`not-json`),
			},
			wantErr: true,
		},
		{
			name: "map missing geometry",
			payload: LocationPayload{
				Source: SourceMap,
			},
			wantErr: true,
		},
		{
			name: "text with description",
			payload: LocationPayload{
				Source:      SourceText,
				Description: "Village near the old bridge",
			},
			wantErr: false,
		},
		{
			name: "text missing description",
			payload: LocationPayload{
				Source: SourceText,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMapValidationAllowsExtraFields(t *testing.T) {
	payload := LocationPayload{
		Source:      SourceMap,
		Name:        "Majuli Island",
		Description: "Northern side near the embankment",
		Latitude:    floatPointer(26.987),
		Longitude:   floatPointer(94.223),
		Geometry:    []byte(`{"type":"Polygon","coordinates":[]}`),
	}

	assert.NoError(t, payload.Validate())
}
