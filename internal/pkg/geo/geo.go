// Package geo provides generic geographic concepts that are independent of any
// database or storage format.
package geo

// Coordinate is a geographic point in WGS84 decimal degrees.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// Valid reports whether the coordinate is within WGS84 bounds.
func (c Coordinate) Valid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 &&
		c.Longitude >= -180 && c.Longitude <= 180
}
