package db

import (
	"context"
	"errors"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
)

var (
	ErrInvalidArg     = errors.New("invalid argument")
	ErrEmptyArg       = errors.New("empty argument not allowed")
	ErrRecordNotFound = errors.New("record not found")
)

//go:generate mockgen -destination=../mocks/mock_db.go -package=mocks github.com/ganeshdipdumbare/scootin-aboot-journey/db DB
// DB interface to interact with database
type DB interface {
	// scooter functions
	GetAvailableScootersWithinRadius(ctx context.Context, location *domain.GeoLocation, radius int) ([]domain.Scooter, error)
	GetScooterByID(ctx context.Context, scooterID string) (*domain.Scooter, error)
	UpdateScooter(ctx context.Context, updatedScooter *domain.Scooter) (*domain.Scooter, error)
	GetAllScooters(ctx context.Context) ([]domain.Scooter, error)
	InsertTripEvent(ctx context.Context, event *domain.TripEvent) error
	GetAllTripEvents(ctx context.Context) ([]domain.TripEvent, error)

	// user functions
	GetAllUsers(ctx context.Context) ([]domain.User, error)

	Disconnect(ctx context.Context) error
}
