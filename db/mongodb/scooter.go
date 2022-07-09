package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/db"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Scooter represents scooter DB record
type Scooter struct {
	InternalID    primitive.ObjectID `bson:"_id,omitempty"`
	ID            string             `bson:"id"`
	Name          string             `bson:"name"`
	Location      GeoLocation        `bson:"location"`
	CurrentUserID *string            `bson:"current_user_id,omitempty"`
	IsAvailable   bool               `bson:"is_available"`
}

// transformToDBScooter creates and returns scooter DB record from domain scooter record
func transformToDBScooter(scooter *domain.Scooter) (*Scooter, error) {
	if scooter == nil {
		return nil, db.ErrInvalidArg
	}

	location := GeoLocation{
		Type: GeoJSONPointType,
		Coordinates: []float64{
			scooter.Location.Latitude,
			scooter.Location.Longitude,
		},
	}

	scooterDB := &Scooter{
		ID:            scooter.ID,
		Name:          scooter.Name,
		Location:      location,
		IsAvailable:   scooter.IsAvailable,
		CurrentUserID: scooter.CurrentUserID,
	}

	return scooterDB, nil
}

// transformToDomainScooter creates and returns scooter domain record from DB scooter record
func transformToDomainScooter(scooter *Scooter) (*domain.Scooter, error) {
	if scooter == nil {
		return nil, db.ErrInvalidArg
	}

	scooterDomain := &domain.Scooter{
		ID:   scooter.ID,
		Name: scooter.Name,
		Location: domain.GeoLocation{
			Latitude:  scooter.Location.Coordinates[0],
			Longitude: scooter.Location.Coordinates[1],
		},
		CurrentUserID: scooter.CurrentUserID,
		IsAvailable:   scooter.IsAvailable,
	}

	return scooterDomain, nil
}

// getScootersByFilter returns all the scooters with given filter
// e.g. get all the scooters which are available using filter: is_available=true
func (m *mongoDetails) getScootersByFilter(ctx context.Context, filter bson.M) ([]domain.Scooter, error) {
	cur, err := m.ScooterCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	scooters := []Scooter{}
	for cur.Next(context.TODO()) {
		var scooter Scooter
		err := cur.Decode(&scooter)
		if err != nil {
			return nil, err
		}

		scooters = append(scooters, scooter)
	}
	result := []domain.Scooter{}
	for _, scooter := range scooters {
		r, err := transformToDomainScooter(&scooter)
		if err != nil {
			return nil, err
		}
		result = append(result, *r)
	}

	return result, nil
}

// GetAvailableScootersWithinRadius returns available scooters which are within
// radius from the location in nearest first sorted order.
func (m *mongoDetails) GetAvailableScootersWithinRadius(ctx context.Context, location *domain.GeoLocation, radius int) ([]domain.Scooter, error) {
	if location == nil {
		return nil, fmt.Errorf("location: %w", db.ErrInvalidArg)
	}

	if radius == 0 {
		return nil, fmt.Errorf("radius: %w", db.ErrInvalidArg)
	}

	filter := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        GeoJSONPointType,
					"coordinates": []float64{location.Latitude, location.Longitude},
				},
				"$maxDistance": radius,
			},
		},
		"is_available": true,
	}

	return m.getScootersByFilter(ctx, filter)
}

// UpdateScooter updates scooter with the given scooter record
func (m *mongoDetails) UpdateScooter(ctx context.Context, scooter *domain.Scooter) (*domain.Scooter, error) {
	if scooter == nil {
		return nil, fmt.Errorf("scooter: %w", db.ErrInvalidArg)
	}

	filter := bson.M{
		"id": scooter.ID,
	}
	dbScooter, err := transformToDBScooter(scooter)
	if err != nil {
		return nil, err
	}

	updateFields := bson.M{
		"$set": bson.M{
			"id":              dbScooter.ID,
			"name":            dbScooter.Name,
			"location":        dbScooter.Location,
			"is_available":    dbScooter.IsAvailable,
			"current_user_id": dbScooter.CurrentUserID,
		},
	}
	_, err = m.ScooterCollection.UpdateOne(ctx, filter, updateFields)
	if err != nil {
		return nil, err
	}
	return scooter, nil
}

// GetAllScooters returns all the scooters in the system
func (m *mongoDetails) GetAllScooters(ctx context.Context) ([]domain.Scooter, error) {
	filter := bson.M{}
	return m.getScootersByFilter(ctx, filter)
}

// GetScooterByID returns scooter for given id, if not found returns error
func (m *mongoDetails) GetScooterByID(ctx context.Context, scooterID string) (*domain.Scooter, error) {
	filter := primitive.M{"id": scooterID}
	var record Scooter
	err := m.ScooterCollection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}
	return transformToDomainScooter(&record)
}
