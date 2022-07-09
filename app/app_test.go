package app

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/db"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AppTestSuite struct {
	suite.Suite
	Database       *mocks.MockDB
	MockController *gomock.Controller
}

// SetupTest runs before every test
func (suite *AppTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.MockController = mockCtrl
	suite.Database = mocks.NewMockDB(mockCtrl)
}

// TearDownTest runs after every test
func (suite *AppTestSuite) TearDownTest() {
	suite.MockController.Finish()
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (suite *AppTestSuite) TestNewApp() {
	t := suite.T()

	type args struct {
		database db.DB
	}
	tests := []struct {
		name    string
		args    args
		want    App
		wantErr bool
	}{
		{
			name: "should return app when valid input db",
			args: args{
				database: suite.Database,
			},
			want: &appDetails{
				database: suite.Database,
			},
			wantErr: false,
		},
		{
			name: "should return error when nil input db",
			args: args{
				database: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApp(tt.args.database)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *AppTestSuite) TestGetNearbyAvailableScooters() {
	t := suite.T()
	database := suite.Database
	ctx := context.Background()

	nearbyScooters := []domain.Scooter{
		{
			ID:   "testid",
			Name: "testscooter",
			Location: domain.GeoLocation{
				Latitude:  0,
				Longitude: 0,
			},
			IsAvailable: true,
		},
	}

	type fields struct {
		database db.DB
	}
	type args struct {
		ctx      context.Context
		location domain.GeoLocation
		radius   int
	}
	tests := []struct {
		name    string
		prepare func()
		fields  fields
		args    args
		want    []domain.Scooter
		wantErr bool
	}{
		{
			name: "should return error for invalid arg",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:      ctx,
				location: domain.GeoLocation{},
				radius:   0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return scooters for valid arg",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:      ctx,
				location: domain.GeoLocation{},
				radius:   10,
			},
			want:    nearbyScooters,
			wantErr: false,
			prepare: func() {
				database.EXPECT().GetAvailableScootersWithinRadius(ctx, gomock.Any(), gomock.Any()).Return(nearbyScooters, nil).Times(1)
			},
		},
		{
			name: "should return error if db returns empty arg err",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:      ctx,
				location: domain.GeoLocation{},
				radius:   10,
			},
			want:    nil,
			wantErr: true,
			prepare: func() {
				database.EXPECT().GetAvailableScootersWithinRadius(ctx, gomock.Any(), gomock.Any()).Return(nil, db.ErrEmptyArg).Times(1)
			},
		},
		{
			name: "should return error if db returns invalid arg err",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:      ctx,
				location: domain.GeoLocation{},
				radius:   10,
			},
			want:    nil,
			wantErr: true,
			prepare: func() {
				database.EXPECT().GetAvailableScootersWithinRadius(ctx, gomock.Any(), gomock.Any()).Return(nil, db.ErrInvalidArg).Times(1)
			},
		},
		{
			name: "should return error if db record not found err",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:      ctx,
				location: domain.GeoLocation{},
				radius:   10,
			},
			want:    nil,
			wantErr: true,
			prepare: func() {
				database.EXPECT().GetAvailableScootersWithinRadius(ctx, gomock.Any(), gomock.Any()).Return(nil, db.ErrRecordNotFound).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &appDetails{
				database: tt.fields.database,
			}

			if tt.prepare != nil {
				tt.prepare()
			}

			got, err := a.GetNearbyAvailableScooters(tt.args.ctx, tt.args.location, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNearbyAvailableScooters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNearbyAvailableScooters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *AppTestSuite) TestBeginTrip() {
	t := suite.T()
	database := suite.Database
	ctx := context.Background()

	type fields struct {
		database db.DB
	}
	type args struct {
		ctx       context.Context
		userID    string
		scooterID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func()
		wantErr bool
	}{
		{
			name: "should return error for empty userID",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "",
				scooterID: "scooterid",
			},
			prepare: func() {},
			wantErr: true,
		},
		{
			name: "should return error for empty scooterID",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "userid",
				scooterID: "",
			},
			prepare: func() {},
			wantErr: true,
		},
		{
			name: "should return error if scooter not found",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "userid",
				scooterID: "scooterid",
			},
			prepare: func() {
				database.EXPECT().GetScooterByID(ctx, gomock.Any()).Return(nil, db.ErrRecordNotFound).Times(1)
			},
			wantErr: true,
		},
		{
			name: "should return error if db error while fetching scooter by id",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "userid",
				scooterID: "scooterid",
			},
			prepare: func() {
				database.EXPECT().GetScooterByID(ctx, gomock.Any()).Return(nil, errors.New("internal error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "should return error if scooter is unavailable",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "userid",
				scooterID: "scooterid",
			},
			prepare: func() {
				database.EXPECT().GetScooterByID(ctx, gomock.Any()).Return(&domain.Scooter{
					ID:            "scooterid",
					Name:          "scooter 1",
					Location:      domain.GeoLocation{},
					CurrentUserID: nil,
					IsAvailable:   false,
				}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name: "should return error if update scooter failed",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "userid",
				scooterID: "scooterid",
			},
			prepare: func() {
				currentScooter := &domain.Scooter{
					ID:            "scooterid",
					Name:          "scooter 1",
					Location:      domain.GeoLocation{},
					CurrentUserID: nil,
					IsAvailable:   true,
				}
				gomock.InOrder(
					database.EXPECT().GetScooterByID(ctx, gomock.Any()).Return(currentScooter, nil).Times(1),
					database.EXPECT().UpdateScooter(ctx, gomock.Any()).Return(nil, errors.New("internal error")).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "should return success if trip is begin successfully",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				userID:    "userid",
				scooterID: "scooterid",
			},
			prepare: func() {
				currentScooter := &domain.Scooter{
					ID:            "scooterid",
					Name:          "scooter 1",
					Location:      domain.GeoLocation{},
					CurrentUserID: nil,
					IsAvailable:   true,
				}

				userID := "userid"
				updatedScooter := *currentScooter
				updatedScooter.CurrentUserID = &userID
				updatedScooter.IsAvailable = false

				gomock.InOrder(
					database.EXPECT().GetScooterByID(ctx, gomock.Any()).Return(currentScooter, nil).Times(1),
					database.EXPECT().UpdateScooter(ctx, &updatedScooter).Return(&updatedScooter, nil).Times(1),
				)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &appDetails{
				database: tt.fields.database,
			}
			tt.prepare()
			if err := a.BeginTrip(tt.args.ctx, tt.args.userID, tt.args.scooterID); (err != nil) != tt.wantErr {
				t.Errorf("appDetails.BeginTrip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
