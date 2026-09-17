package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinateValid(t *testing.T) {
	assert.True(t, Coordinate{Latitude: 0, Longitude: 0}.Valid())
	assert.True(t, Coordinate{Latitude: 26.18, Longitude: 91.73}.Valid())
	assert.True(t, Coordinate{Latitude: 90, Longitude: 180}.Valid())
	assert.True(t, Coordinate{Latitude: -90, Longitude: -180}.Valid())

	assert.False(t, Coordinate{Latitude: 90.1, Longitude: 0}.Valid())
	assert.False(t, Coordinate{Latitude: -90.1, Longitude: 0}.Valid())
	assert.False(t, Coordinate{Latitude: 0, Longitude: 180.1}.Valid())
	assert.False(t, Coordinate{Latitude: 0, Longitude: -180.1}.Valid())
}
