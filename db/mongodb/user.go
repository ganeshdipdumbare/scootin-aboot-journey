package mongodb

import (
	"context"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/db"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents user DB record
type User struct {
	InternalID primitive.ObjectID `bson:"_id,omitempty"`
	ID         string             `bson:"id"`
	Name       string             `bson:"name"`
}

// transformToDomainUser convert db user to domain user
func transformToDomainUser(user *User) (*domain.User, error) {
	if user == nil {
		return nil, db.ErrInvalidArg
	}

	domainUser := &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}
	return domainUser, nil
}

// GetAllUsers returns all the users
func (m *mongoDetails) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	filter := bson.M{}
	cur, err := m.UserCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for cur.Next(context.TODO()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	result := []domain.User{}
	for _, user := range users {
		r, err := transformToDomainUser(&user)
		if err != nil {
			return nil, err
		}
		result = append(result, *r)
	}

	return result, nil
}
