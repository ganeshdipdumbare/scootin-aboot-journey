package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/app"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	App            *mocks.MockApp
	MockController *gomock.Controller
}

// SetupTest runs before every test
func (suite *HandlerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.MockController = mockCtrl
	suite.App = mocks.NewMockApp(mockCtrl)
}

// TearDownTest runs after every test
func (suite *HandlerTestSuite) TearDownTest() {
	suite.MockController.Finish()
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) Test_getAvailableScooters() {
	t := suite.T()
	appInstance := suite.App
	api := &apiDetails{
		app:    appInstance,
		apiKey: "testkey",
	}
	router := api.setupRouter()
	availableScooterApiPath := "/api/v1/auth/user/available-scooters"

	type args struct {
		url string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    want
	}{
		{
			name:    "should return error invalid api key",
			prepare: func() {},
			args: args{
				url: availableScooterApiPath + "?api_key=invalid",
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:    "should return error for invalid latitude",
			prepare: func() {},
			args: args{
				url: availableScooterApiPath + "?latitude=invalid&api_key=testkey",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "should return error for invalid longitude",
			prepare: func() {},
			args: args{
				url: availableScooterApiPath + "?longitude=invalid&latitude=0.0&api_key=testkey",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "should return error for invalid radius",
			prepare: func() {},
			args: args{
				url: availableScooterApiPath + "?longitude=0.0&latitude=0.0&radius=invalid&api_key=testkey",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return error if get availabe scooters returns error",
			prepare: func() {
				appInstance.EXPECT().GetNearbyAvailableScooters(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, app.ErrEmptyArg).Times(1)
			},
			args: args{
				url: availableScooterApiPath + "?longitude=0.0&latitude=0.0&radius=2&api_key=testkey",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return scooters for valid args",
			prepare: func() {
				respScooter := domain.Scooter{
					ID:   "testscooter",
					Name: "testname",
					Location: domain.GeoLocation{
						Latitude:  0.0,
						Longitude: 0.0,
					},
					CurrentUserID: nil,
					IsAvailable:   true,
				}
				appInstance.EXPECT().GetNearbyAvailableScooters(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Scooter{respScooter}, nil).Times(1)
			},
			args: args{
				url: availableScooterApiPath + "?longitude=0.0&latitude=0.0&radius=2&api_key=testkey",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.args.url, nil)
			router.ServeHTTP(w, req)

			if tt.want.statusCode != w.Code {
				t.Errorf("getAvailableScooters() status code  = %v, want status code %v", w.Code, tt.want.statusCode)
				return
			}
		})
	}
}
