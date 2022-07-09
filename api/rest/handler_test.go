package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (suite *HandlerTestSuite) Test_beginTrip() {
	t := suite.T()
	appInstance := suite.App
	api := &apiDetails{
		app:    appInstance,
		apiKey: "testkey",
	}
	router := api.setupRouter()
	beginTripApiPath := "/api/v1/auth/user/begin-trip"

	type args struct {
		url  string
		body io.Reader
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
			name:    "should return error for invalid api key",
			prepare: func() {},
			args: args{
				url: beginTripApiPath + "?api_key=invalid",
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:    "should return error for invalid body param",
			prepare: func() {},
			args: args{
				url: beginTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":"invalidid",
					"user_id":"invalidid"
				}`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "should return error for invalid body param type",
			prepare: func() {},
			args: args{
				url: beginTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":1,
					"user_id":"invalidid"
				}`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return error if app BeginTrip returns error",
			prepare: func() {
				appInstance.EXPECT().BeginTrip(gomock.Any(), gomock.Any(), gomock.Any()).Return(app.ErrRecordNotFound).Times(1)
			},
			args: args{
				url: beginTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":"f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"user_id":"f3b9842c-182a-418b-92fd-95d4f46414c5"
				}`),
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "should return success if app BeginTrip returns success",
			prepare: func() {
				appInstance.EXPECT().BeginTrip(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				url: beginTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":"f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"user_id":"f3b9842c-182a-418b-92fd-95d4f46414c5"
				}`),
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
			req, _ := http.NewRequest(http.MethodPut, tt.args.url, tt.args.body)
			router.ServeHTTP(w, req)

			if tt.want.statusCode != w.Code {
				t.Errorf("beginTrip() status code  = %v, want status code %v", w.Code, tt.want.statusCode)
				return
			}
		})
	}
}

func (suite *HandlerTestSuite) Test_endTrip() {
	t := suite.T()
	appInstance := suite.App
	api := &apiDetails{
		app:    appInstance,
		apiKey: "testkey",
	}
	router := api.setupRouter()
	endTripApiPath := "/api/v1/auth/user/end-trip"

	type args struct {
		url  string
		body io.Reader
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
			name:    "should return error for invalid api key",
			prepare: func() {},
			args: args{
				url: endTripApiPath + "?api_key=invalid",
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:    "should return error for invalid body param",
			prepare: func() {},
			args: args{
				url: endTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":"invalidid",
					"user_id":"invalidid"
				}`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "should return error for invalid body param type",
			prepare: func() {},
			args: args{
				url: endTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":1,
					"user_id":"invalidid"
				}`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return error if app EndTrip returns error",
			prepare: func() {
				appInstance.EXPECT().EndTrip(gomock.Any(), gomock.Any(), gomock.Any()).Return(app.ErrRecordNotFound).Times(1)
			},
			args: args{
				url: endTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":"f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"user_id":"f3b9842c-182a-418b-92fd-95d4f46414c5"
				}`),
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "should return success if app EndTrip returns success",
			prepare: func() {
				appInstance.EXPECT().EndTrip(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				url: endTripApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"scooter_id":"f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"user_id":"f3b9842c-182a-418b-92fd-95d4f46414c5"
				}`),
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
			req, _ := http.NewRequest(http.MethodPut, tt.args.url, tt.args.body)
			router.ServeHTTP(w, req)

			if tt.want.statusCode != w.Code {
				t.Errorf("endTrip() status code  = %v, want status code %v", w.Code, tt.want.statusCode)
				return
			}
		})
	}
}

func (suite *HandlerTestSuite) Test_saveScooterTripEvent() {
	t := suite.T()
	appInstance := suite.App
	api := &apiDetails{
		app:    appInstance,
		apiKey: "testkey",
	}
	router := api.setupRouter()
	saveTripEventApiPath := "/api/v1/auth/scooter/trip-event"

	type args struct {
		url  string
		body io.Reader
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
			name:    "should return error for invalid api key",
			prepare: func() {},
			args: args{
				url: saveTripEventApiPath + "?api_key=invalid",
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:    "should return error for invalid body param",
			prepare: func() {},
			args: args{
				url: saveTripEventApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"created_at": "2022-07-09T17:49:09+00:00",
					"location": {
					  "latitude": 0,
					  "longitude": 0
					},
					"scooter_id": "invalidid",
					"type": "trip_stop",
					"user_id": "f3b9842c-182a-418b-92fd-95d4f46414c5"
				  }`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "should return error for invalid body param type",
			prepare: func() {},
			args: args{
				url: saveTripEventApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"created_at": "2022-07-09T17:49:09+00:00",
					"location": {
					  "latitude": 0,
					  "longitude": 0
					},
					"scooter_id": 0,
					"type": "trip_stop",
					"user_id": "f3b9842c-182a-418b-92fd-95d4f46414c5"
				  }`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "should return error for invalid event type",
			prepare: func() {},
			args: args{
				url: saveTripEventApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"created_at": "2022-07-09T17:49:09+00:00",
					"location": {
					  "latitude": 0,
					  "longitude": 0
					},
					"scooter_id": "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"type": "invalid",
					"user_id": "f3b9842c-182a-418b-92fd-95d4f46414c5"
				  }`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return error if error while saving trip event",
			prepare: func() {
				appInstance.EXPECT().SaveScooterTripEvent(gomock.Any(), gomock.Any()).Return(app.ErrInvalidArg).Times(1)
			},
			args: args{
				url: saveTripEventApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"created_at": "2022-07-09T17:49:09+00:00",
					"location": {
					  "latitude": 0,
					  "longitude": 0
					},
					"scooter_id": "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"type": "trip_start",
					"user_id": "f3b9842c-182a-418b-92fd-95d4f46414c5"
				  }`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return success if the trip event is saved successfully",
			prepare: func() {
				appInstance.EXPECT().SaveScooterTripEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				url: saveTripEventApiPath + "?api_key=testkey",
				body: strings.NewReader(`{
					"created_at": "2022-07-09T17:49:09+00:00",
					"location": {
					  "latitude": 0,
					  "longitude": 0
					},
					"scooter_id": "f691fd32-9b3f-4d71-b9b7-c48213bfd232",
					"type": "trip_start",
					"user_id": "f3b9842c-182a-418b-92fd-95d4f46414c5"
				  }`),
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, tt.args.url, tt.args.body)
			router.ServeHTTP(w, req)

			if tt.want.statusCode != w.Code {
				t.Errorf("saveScooterTripEvent() status code  = %v, want status code %v", w.Code, tt.want.statusCode)
				return
			}
		})
	}
}
