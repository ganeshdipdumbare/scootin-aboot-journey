package mongodb

import (
	"context"
	"reflect"
	"testing"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

func (suite *MongoTestSuite) Test_mongoDetails_GetAllUsers() {
	mgoC := suite.TestContainer
	t := suite.T()
	dbName := "testdb"

	client, err := connectAndMigrateTestData(mgoC, dbName)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		client              *mongo.Client
		dbName              string
		ScooterCollection   *mongo.Collection
		UserCollection      *mongo.Collection
		TripEventCollection *mongo.Collection
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		wantErr bool
	}{
		{
			name: "should return success",
			fields: fields{
				client:         client,
				dbName:         dbName,
				UserCollection: client.Database(dbName).Collection(userCollectionName),
			},
			args: args{
				ctx: context.Background(),
			},
			want: []domain.User{
				{
					ID:   "f3b9842c-182a-418b-92fd-95d4f46414c5",
					Name: "User 1",
				},
				{
					ID:   "6124edb7-5099-4147-87e6-0c9b93cd1fdb",
					Name: "User 2",
				},
				{
					ID:   "4668a2f7-c498-4e49-a82e-380c1ede0685",
					Name: "User 3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mongoDetails{
				client:              tt.fields.client,
				dbName:              tt.fields.dbName,
				ScooterCollection:   tt.fields.ScooterCollection,
				UserCollection:      tt.fields.UserCollection,
				TripEventCollection: tt.fields.TripEventCollection,
			}
			got, err := m.GetAllUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
