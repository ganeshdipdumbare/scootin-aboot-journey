package domain

// Scooter represents scooter details
type Scooter struct {
	ID            string
	Name          string
	Location      GeoLocation
	CurrentUserID *string
	IsAvailable   bool
}
