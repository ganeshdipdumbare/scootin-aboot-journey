package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

func (suite *MongoTestSuite) TestInsertTripEvent() {
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
		ctx       context.Context
		tripEvent *domain.TripEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should return error for nil arg",
			fields: fields{
				client:              client,
				dbName:              dbName,
				TripEventCollection: client.Database(dbName).Collection(tripEventCollectionName),
			},
			args: args{
				ctx:       context.Background(),
				tripEvent: nil,
			},
			wantErr: true,
		},
		{
			name: "should return success for valid input arg",
			fields: fields{
				client:              client,
				dbName:              dbName,
				TripEventCollection: client.Database(dbName).Collection(tripEventCollectionName),
			},
			args: args{
				ctx: context.Background(),
				tripEvent: &domain.TripEvent{
					UserID:    "userid",
					ScooterID: "scooterid",
					Location: domain.GeoLocation{
						Latitude:  0.0,
						Longitude: 0.0,
					},
					Type:      domain.TripStartEvent,
					CreatedAt: time.Now(),
				},
			},
			wantErr: false,
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
			if err := m.InsertTripEvent(tt.args.ctx, tt.args.tripEvent); (err != nil) != tt.wantErr {
				t.Errorf("InsertTripEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *MongoTestSuite) TestGetAllTripEvents() {
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
		wantErr bool
	}{
		{
			name: "should return success",
			fields: fields{
				client:              client,
				dbName:              dbName,
				TripEventCollection: client.Database(dbName).Collection(tripEventCollectionName),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
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
			_, err := m.GetAllTripEvents(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTripEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
