package mongodb

type GeoJSONType string

const (
	GeoJSONPointType GeoJSONType = "Point"
)

// GeoLocation represents the location in mongodb
// for more details, visit - https://www.mongodb.com/docs/manual/geospatial-queries/
type GeoLocation struct {
	Type        GeoJSONType `json:"type" bson:"type"`
	Coordinates []float64   `json:"coordinates" bson:"coordinates"`
}
