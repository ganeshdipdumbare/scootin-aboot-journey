package testclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"github.com/go-resty/resty/v2"
)

var (
	errInvalidRespStatusCode = errors.New("invalid status code in response")
)

type getAvailableScootersResponse struct {
	Scooters []scooter `json:"scooters"`
}

type scooter struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Location      geoLocation `json:"location"`
	CurrentUserID *string     `json:"current_user_id"`
	IsAvailable   bool        `json:"is_available"`
}

type geoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type beginTripRequest struct {
	UserID    string `json:"user_id"`
	ScooterID string `json:"scooter_id"`
}

type beginTripResponse struct {
	UserID    string `json:"user_id"`
	ScooterID string `json:"scooter_id"`
}

type endTripRequest struct {
	UserID    string      `json:"user_id"`
	ScooterID string      `json:"scooter_id"`
	Location  geoLocation `json:"location"`
}

type endTripResponse struct {
	UserID    string      `json:"user_id"`
	ScooterID string      `json:"scooter_id"`
	Location  geoLocation `json:"location"`
}

type saveScooterTripEventRequest struct {
	UserID    string      `json:"user_id"`
	ScooterID string      `json:"scooter_id"`
	Location  geoLocation `json:"location"`
	CreatedAt time.Time   `json:"created_at"`
	Type      string      `json:"type"`
}

type saveScooterTripEventResponse struct {
	Success bool `json:"success"`
}

type testClient struct {
	userID          string
	currentLocation *domain.GeoLocation
	radius          int
	apiKey          string
	travelTime      time.Duration
	restTime        time.Duration
	httpClient      *resty.Client
}

// NewTestClientReq
type NewTestClientReq struct {
	ApiKey          string
	Port            string
	UserID          string
	CurrentLocation *domain.GeoLocation
	Radius          int
	TravelTime      time.Duration
	RestTime        time.Duration
}

func NewTestClient(req *NewTestClientReq) *testClient {
	baseURL := fmt.Sprintf("http://localhost:%s/api/v1", req.Port)
	restyClient := resty.New()
	restyClient = restyClient.SetBaseURL(baseURL)
	return &testClient{
		userID:          req.UserID,
		currentLocation: req.CurrentLocation,
		travelTime:      req.TravelTime,
		restTime:        req.RestTime,
		httpClient:      restyClient,
		radius:          req.Radius,
		apiKey:          req.ApiKey,
	}
}

func (tc *testClient) StartJourney() {
	for {
		scooterID, err := tc.getAvailableScooter()
		if err != nil {
			log.Println(err)
			return
		}

		err = tc.beginTrip(scooterID)
		if err != nil {
			log.Println(err)
			return
		}

		currentLocation := geoLocation{
			Latitude:  tc.currentLocation.Latitude,
			Longitude: tc.currentLocation.Longitude,
		}
		err = tc.saveTripEvent(scooterID, "trip_start", currentLocation)
		if err != nil {
			log.Println(err)
		}

		tc.updateLocationDuringTrip(scooterID)

		err = tc.endTrip(scooterID, tc.currentLocation)
		if err != nil {
			log.Println(err)
			return
		}

		currentLocation = geoLocation{
			Latitude:  tc.currentLocation.Latitude,
			Longitude: tc.currentLocation.Longitude,
		}
		err = tc.saveTripEvent(scooterID, "trip_stop", currentLocation)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(tc.restTime)
	}
}

func (tc *testClient) updateLocationDuringTrip(scooterID string) {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(tc.travelTime)
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("trip over!")
			return
		case <-ticker.C:
			tc.currentLocation = travelTenMeterNorth(tc.currentLocation)
			currentLocation := geoLocation{
				Latitude:  tc.currentLocation.Latitude,
				Longitude: tc.currentLocation.Longitude,
			}
			err := tc.saveTripEvent(scooterID, "trip_location_update", currentLocation)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func travelTenMeterNorth(currentLocation *domain.GeoLocation) *domain.GeoLocation {
	// code to find out the next location in North direction
	r_earth := 6378.0
	currentLocation.Latitude = currentLocation.Latitude + (0.01/r_earth)*(180/math.Pi)
	return currentLocation
}

// getAvailableScooter returns closest scooter from user's current location
func (tc *testClient) getAvailableScooter() (string, error) {
	resp, err := tc.httpClient.R().
		SetQueryParams(map[string]string{
			"latitude":  fmt.Sprintf("%f", tc.currentLocation.Latitude),
			"longitude": fmt.Sprintf("%f", tc.currentLocation.Longitude),
			"radius":    fmt.Sprintf("%d", tc.radius),
			"api_key":   tc.apiKey,
		}).
		SetHeader("Accept", "application/json").
		Get("/auth/user/available-scooters")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errInvalidRespStatusCode
	}

	scootersResp := getAvailableScootersResponse{}
	bs := resp.Body()
	err = json.Unmarshal(bs, &scootersResp)
	if err != nil {
		return "", err
	}

	if len(scootersResp.Scooters) == 0 {
		return "", errors.New("no scooter available")
	}
	return scootersResp.Scooters[0].ID, nil
}

// beginTrip begin the trip with given scooter id
func (tc *testClient) beginTrip(scooterID string) error {
	beginTripReqBody := beginTripRequest{
		UserID:    tc.userID,
		ScooterID: scooterID,
	}

	resp, err := tc.httpClient.R().
		SetQueryParams(map[string]string{
			"api_key": tc.apiKey,
		}).SetBody(beginTripReqBody).
		SetHeader("Accept", "application/json").
		Put("/auth/user/begin-trip")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errInvalidRespStatusCode
	}
	return nil
}

// beginTrip end the trip with given scooter id
func (tc *testClient) endTrip(scooterID string, location *domain.GeoLocation) error {
	endTripReqBody := endTripRequest{
		UserID:    tc.userID,
		ScooterID: scooterID,
		Location: geoLocation{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
	}

	resp, err := tc.httpClient.R().
		SetQueryParams(map[string]string{
			"api_key": tc.apiKey,
		}).SetBody(endTripReqBody).
		SetHeader("Accept", "application/json").
		Put("/auth/user/end-trip")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errInvalidRespStatusCode
	}
	return nil
}

func (tc *testClient) saveTripEvent(scooterID, eventType string, location geoLocation) error {
	saveTripEventReqBody := saveScooterTripEventRequest{
		UserID:    tc.userID,
		ScooterID: scooterID,
		Location:  location,
		CreatedAt: time.Now().UTC(),
		Type:      eventType,
	}

	resp, err := tc.httpClient.R().
		SetQueryParams(map[string]string{
			"api_key": tc.apiKey,
		}).SetBody(saveTripEventReqBody).
		SetHeader("Accept", "application/json").
		Post("/auth/scooter/trip-event")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return errInvalidRespStatusCode
	}

	return nil
}
