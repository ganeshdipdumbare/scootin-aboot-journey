package mongodb

import (
	"context"
	"reflect"
	"testing"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

func (suite *MongoTestSuite) TestGetScootersWithinRadius() {
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
		ScooterCollectoin   *mongo.Collection
		UserCollection      *mongo.Collection
		TripEventCollection *mongo.Collection
	}
	type args struct {
		ctx      context.Context
		location *domain.GeoLocation
		radius   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Scooter
		wantErr bool
	}{
		{
			name: "should return empty list for the location not near",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollectoin: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx: context.Background(),
				location: &domain.GeoLocation{
					Latitude:  0.0,
					Longitude: 0.0,
				},
				radius: 10.0,
			},
			want:    []domain.Scooter{},
			wantErr: false,
		},
		{
			name: "should return empty list for the location near scooter 1",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollectoin: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx: context.Background(),
				location: &domain.GeoLocation{
					Latitude:  -73.856077,
					Longitude: 40.848447,
				},
				radius: 10.0,
			},
			want: []domain.Scooter{
				{
					ID:   "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					Name: "Scooter 1",
					Location: domain.GeoLocation{
						Latitude:  -73.856077,
						Longitude: 40.848447,
					},
					IsAvailable: true,
				},
			},
			wantErr: false,
		},
		{
			name: "should return error nil location",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollectoin: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx: context.Background(),
				location: &domain.GeoLocation{
					Latitude:  0.0,
					Longitude: 0.0,
				},
				radius: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error for 0 radius",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollectoin: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx:      context.Background(),
				location: nil,
				radius:   10.0,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mongoDetails{
				client:              tt.fields.client,
				dbName:              tt.fields.dbName,
				ScooterCollection:   tt.fields.ScooterCollectoin,
				UserCollection:      tt.fields.UserCollection,
				TripEventCollection: tt.fields.TripEventCollection,
			}
			got, err := m.GetAvailableScootersWithinRadius(tt.args.ctx, tt.args.location, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvailableScootersWithinRadius() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAvailableScootersWithinRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *MongoTestSuite) TestUpdateScooter() {
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
		ctx     context.Context
		scooter *domain.Scooter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Scooter
		wantErr bool
	}{
		{
			name: "should return error for nil scooter arg",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollection: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx:     context.Background(),
				scooter: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return success for valid scooter arg",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollection: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx: context.Background(),
				scooter: &domain.Scooter{
					ID:   "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					Name: "Scooter 1",
					Location: domain.GeoLocation{
						Latitude:  -73.856077,
						Longitude: 40.848447,
					},
					IsAvailable: false,
				},
			},
			want: &domain.Scooter{
				ID:   "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
				Name: "Scooter 1",
				Location: domain.GeoLocation{
					Latitude:  -73.856077,
					Longitude: 40.848447,
				},
				IsAvailable: false,
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
			got, err := m.UpdateScooter(tt.args.ctx, tt.args.scooter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateScooter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateScooter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *MongoTestSuite) TestGetAllScooters() {
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
		want    []domain.Scooter
		wantErr bool
	}{
		{

			name: "should return success for valid scooter arg",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollection: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx: context.Background(),
			},
			want: []domain.Scooter{
				{
					ID: "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					Location: domain.GeoLocation{
						Latitude:  -73.856077,
						Longitude: 40.848447,
					},
					Name:        "Scooter 1",
					IsAvailable: true,
				},
				{

					ID: "10f8cfb7-7764-4b75-acca-cc17d2b07d59",
					Location: domain.GeoLocation{
						Latitude:  -73.961704,
						Longitude: 40.662942,
					},
					Name:        "Scooter 2",
					IsAvailable: true,
				},
				{

					ID: "9360f883-cf55-421e-b21a-1752167f5221",
					Location: domain.GeoLocation{
						Latitude:  -73.98241999999999,
						Longitude: 40.579505,
					},
					Name:        "Scooter 3",
					IsAvailable: true,
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
			got, err := m.GetAllScooters(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllScooters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllScooters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *MongoTestSuite) TestGetScooterByID() {
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
		scooterID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Scooter
		wantErr bool
	}{
		{
			name: "should return error if scooter not for id",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollection: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx:       context.Background(),
				scooterID: "invalidscooter",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return success if scooter found for id",
			fields: fields{
				client:            client,
				dbName:            dbName,
				ScooterCollection: client.Database(dbName).Collection(scooterCollectionName),
			},
			args: args{
				ctx:       context.Background(),
				scooterID: "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
			},
			want: &domain.Scooter{
				ID: "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
				Location: domain.GeoLocation{
					Latitude:  -73.856077,
					Longitude: 40.848447,
				},
				Name:        "Scooter 1",
				IsAvailable: true,
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
			got, err := m.GetScooterByID(tt.args.ctx, tt.args.scooterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoDetails.GetScooterByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoDetails.GetScooterByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
