package database

// GeometryToText converts a PostGIS geometry value returned by sqlc (string or
// []byte) into its text form (e.g. the GeoJSON produced by ST_AsGeoJSON). An
// empty or missing value yields "".
func GeometryToText(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}
