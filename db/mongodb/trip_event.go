package mongodb

import (
	"context"
	"time"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/db"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TripEvent represents trip event DB record
type TripEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	ScooterID string             `bson:"scooter_id"`
	Location  GeoLocation        `bson:"location"`
	Type      string             `bson:"type"`
	CreatedAt time.Time          `bson:"created_at"`
}

// transformToDBTripEvent creates db trip event record from domain record
func transformToDBTripEvent(tripEvent *domain.TripEvent) (*TripEvent, error) {
	if tripEvent == nil {
		return nil, db.ErrInvalidArg
	}

	location := GeoLocation{
		Type:        GeoJSONPointType,
		Coordinates: []float64{tripEvent.Location.Latitude, tripEvent.Location.Longitude},
	}

	dbTripEvent := &TripEvent{
		ID:        primitive.NewObjectID(),
		UserID:    tripEvent.UserID,
		ScooterID: tripEvent.ScooterID,
		Location:  location,
		Type:      string(tripEvent.Type),
		CreatedAt: tripEvent.CreatedAt,
	}
	return dbTripEvent, nil
}

// transformToDomainTripEvent creates domain trip event record from db record
func transformToDomainTripEvent(tripEvent *TripEvent) (*domain.TripEvent, error) {
	if tripEvent == nil {
		return nil, db.ErrInvalidArg
	}

	location := domain.GeoLocation{
		Latitude:  tripEvent.Location.Coordinates[0],
		Longitude: tripEvent.Location.Coordinates[1],
	}

	domainTripEvent := &domain.TripEvent{
		ID:        tripEvent.ID.Hex(),
		UserID:    tripEvent.UserID,
		ScooterID: tripEvent.ScooterID,
		Location:  location,
		Type:      domain.TripEventType(tripEvent.Type),
		CreatedAt: tripEvent.CreatedAt,
	}
	return domainTripEvent, nil
}

// InsertTripEvent inserts trip event in the trip_event collection
func (m *mongoDetails) InsertTripEvent(ctx context.Context, tripEvent *domain.TripEvent) error {
	dbTripEvent, err := transformToDBTripEvent(tripEvent)
	if err != nil {
		return err
	}

	_, err = m.TripEventCollection.InsertOne(ctx, dbTripEvent)
	return err
}

// GetAllTripEvents get all trip event
func (m *mongoDetails) GetAllTripEvents(ctx context.Context) ([]domain.TripEvent, error) {
	filter := bson.M{}
	cur, err := m.TripEventCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	tripEvents := []TripEvent{}
	for cur.Next(context.TODO()) {
		var tripEvent TripEvent
		err := cur.Decode(&tripEvent)
		if err != nil {
			return nil, err
		}

		tripEvents = append(tripEvents, tripEvent)
	}
	result := []domain.TripEvent{}
	for _, scooter := range tripEvents {
		r, err := transformToDomainTripEvent(&scooter)
		if err != nil {
			return nil, err
		}
		result = append(result, *r)
	}

	return result, nil
}
