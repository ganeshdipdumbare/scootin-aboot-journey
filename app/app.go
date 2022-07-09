package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/db"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
)

var (
	ErrInvalidArg          = errors.New("invalid argument")
	ErrEmptyArg            = errors.New("empty argument")
	ErrRecordNotFound      = errors.New("record not found")
	ErrOperationNotAllowed = errors.New("operation not allowed")
)

//go:generate mockgen -destination=../mocks/mock_app.go -package=mocks github.com/ganeshdipdumbare/scootin-aboot-journey/app App
// App interface which consists of business logic/use cases
type App interface {
	GetNearbyAvailableScooters(ctx context.Context, location domain.GeoLocation, radius int) ([]domain.Scooter, error)
	BeginTrip(ctx context.Context, userID string, scooterID string) error
}

type appDetails struct {
	database db.DB
}

// NewApp creates new app instance
func NewApp(database db.DB) (App, error) {
	if database == nil {
		return nil, fmt.Errorf("database: %w", ErrInvalidArg)
	}

	return &appDetails{
		database: database,
	}, nil
}

// GetNearbyAvailableScooters returns nearby scooters within radius(meters) from
// the location in nearest first sorted order
func (a *appDetails) GetNearbyAvailableScooters(ctx context.Context, location domain.GeoLocation, radius int) ([]domain.Scooter, error) {
	if radius == 0 {
		return nil, ErrInvalidArg
	}

	userLocation := location
	scooters, err := a.database.GetAvailableScootersWithinRadius(ctx, &userLocation, radius)
	if err != nil {
		var returnErr error
		switch {
		case errors.Is(err, db.ErrInvalidArg):
			returnErr = ErrInvalidArg
		case errors.Is(err, db.ErrEmptyArg):
			returnErr = ErrEmptyArg
		default:
			returnErr = err
		}
		return nil, fmt.Errorf("db error while getting scooters: %w", returnErr)
	}

	return scooters, nil
}

// BeginTrip starts trip for given user with given scooter
// scooter record is updated with current user and set to unavailable
// returns error if scooter is not available
func (a *appDetails) BeginTrip(ctx context.Context, userID string, scooterID string) error {
	if userID == "" {
		return fmt.Errorf("userID: %w", ErrEmptyArg)
	}

	if scooterID == "" {
		return fmt.Errorf("scooterID: %w", ErrEmptyArg)
	}

	scooter, err := a.database.GetScooterByID(ctx, scooterID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return fmt.Errorf("scooter not found: %w", ErrRecordNotFound)
		}
		return fmt.Errorf("unable to get scooter: %w", err)
	}

	if !scooter.IsAvailable {
		return fmt.Errorf("scooter is unavailable: %w", ErrOperationNotAllowed)
	}

	updatedScooter := *scooter
	currentUserID := userID
	updatedScooter.CurrentUserID = &currentUserID
	updatedScooter.IsAvailable = false

	_, err = a.database.UpdateScooter(ctx, &updatedScooter)
	if err != nil {
		return fmt.Errorf("unable to update scooter: %w", err)
	}

	return nil
}
