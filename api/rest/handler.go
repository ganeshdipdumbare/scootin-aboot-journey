package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/app"
	docs "github.com/ganeshdipdumbare/scootin-aboot-journey/docs"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	validate *validator.Validate
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
	UserID    string `json:"user_id" validate:"required,uuid4"`
	ScooterID string `json:"scooter_id" validate:"required,uuid4"`
}

type beginTripResponse struct {
	UserID    string `json:"user_id"`
	ScooterID string `json:"scooter_id"`
}

type endTripRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid4"`
	ScooterID string `json:"scooter_id" validate:"required,uuid4"`
}

type endTripResponse struct {
	UserID    string `json:"user_id"`
	ScooterID string `json:"scooter_id"`
}

type saveScooterTripEventRequest struct {
	UserID    string      `json:"user_id" validate:"required,uuid4"`
	ScooterID string      `json:"scooter_id" validate:"required,uuid4"`
	Location  geoLocation `json:"location" validate:"required"`
	CreatedAt time.Time   `json:"created_at" validate:"required"`
	Type      string      `json:"type" validate:"required"`
}

type saveScooterTripEventResponse struct {
	Success bool `json:"success"`
}

type errorRespose struct {
	ErrorMessage string `json:"errorMessage"`
}

func getErrHTTPStatusCode(err error) int {
	httpCode := http.StatusInternalServerError
	switch {
	case errors.Is(err, app.ErrEmptyArg) || errors.Is(err, app.ErrInvalidArg) || errors.Is(err, app.ErrOperationNotAllowed):
		httpCode = http.StatusBadRequest
	case errors.Is(err, app.ErrRecordNotFound):
		httpCode = http.StatusNotFound
	}
	return httpCode
}

func createErrorResponse(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, &errorRespose{
		ErrorMessage: message,
	})
}

func (api *apiDetails) authenticate(c *gin.Context) {
	apiKey := c.Query("api_key")
	fmt.Println("middleware: ", apiKey, api.apiKey)
	if apiKey != api.apiKey {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (api *apiDetails) setupRouter() *gin.Engine {
	validate = validator.New()

	apiV1 := "/api/v1"
	docs.SwaggerInfo.BasePath = apiV1

	r := gin.Default()
	v1group := r.Group(apiV1)
	v1group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authUserGroup := v1group.Group("/auth/user")
	authUserGroup.Use(api.authenticate)
	authUserGroup.GET("/available-scooters", api.getAvailableScooters)
	authUserGroup.PUT("/begin-trip", api.beginTrip)
	authUserGroup.PUT("/end-trip", api.endTrip)

	authScooterGroup := v1group.Group("/auth/scooter")
	authScooterGroup.Use(api.authenticate)
	authScooterGroup.POST("/trip-event", api.saveScooterTripEvent)

	return r
}

// @BasePath /api/v1

// getAvailableScooters godoc
// @Summary returns available scooters within given area
// @Description returns available scooters within given radius sorted by nearest first
// @Tags user-api
// @Accept  json
// @Produce  json
// @Param latitude query number true "latitude"
// @Param longitude query number true "longitude"
// @Param radius query integer true "radius"
// @Param api_key query string true "api_key"
// @Success 200 {object} rest.getAvailableScootersResponse
// @Failure 404 {object} rest.errorRespose
// @Failure 400 {object} rest.errorRespose
// @Failure 500 {object} rest.errorRespose
// @Router /auth/user/available-scooters [get]
func (api *apiDetails) getAvailableScooters(c *gin.Context) {
	lat := c.Query("latitude")
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, "invalid latitude")
		return
	}

	lng := c.Query("longitude")
	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, "invalid longitude")
		return
	}

	rad := c.Query("radius")
	radius, err := strconv.ParseInt(rad, 10, 64)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, "invalid radius")
		return
	}

	userLocation := domain.GeoLocation{
		Latitude:  latitude,
		Longitude: longitude,
	}
	scooters, err := api.app.GetNearbyAvailableScooters(c, userLocation, int(radius))
	if err != nil {
		errStatusCode := getErrHTTPStatusCode(err)
		createErrorResponse(c, errStatusCode, err.Error())
		return
	}

	resp := getAvailableScootersResponse{
		Scooters: []scooter{},
	}
	for _, s := range scooters {
		location := geoLocation{
			Latitude:  s.Location.Latitude,
			Longitude: s.Location.Longitude,
		}

		scooter := scooter{
			ID:            s.ID,
			Name:          s.Name,
			Location:      location,
			CurrentUserID: s.CurrentUserID,
			IsAvailable:   s.IsAvailable,
		}
		resp.Scooters = append(resp.Scooters, scooter)
	}

	c.IndentedJSON(http.StatusOK, resp)
	c.Done()
}

// beginTrip godoc
// @Summary begins the trip
// @Description begins the trip for given user with given scooter, scooter becomes unavailable for other users once the trip begins
// @Tags user-api
// @Accept  json
// @Produce  json
// @Param beginTripRequest body rest.beginTripRequest true "begin trip request"
// @Param api_key query string true "api_key"
// @Success 200 {object} rest.beginTripResponse
// @Failure 404 {object} rest.errorRespose
// @Failure 400 {object} rest.errorRespose
// @Failure 500 {object} rest.errorRespose
// @Router /auth/user/begin-trip [put]
func (api *apiDetails) beginTrip(c *gin.Context) {
	req := &beginTripRequest{}
	err := c.BindJSON(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = api.app.BeginTrip(c, req.UserID, req.ScooterID)
	if err != nil {
		errStatusCode := getErrHTTPStatusCode(err)
		createErrorResponse(c, errStatusCode, err.Error())
		return
	}

	resp := beginTripResponse{
		UserID:    req.UserID,
		ScooterID: req.ScooterID,
	}

	c.IndentedJSON(http.StatusOK, resp)
	c.Done()
}

// endTrip godoc
// @Summary ends the trip
// @Description ends the trip for given user with given scooter, scooter becomes available for other users once the trip ends
// @Tags user-api
// @Accept  json
// @Produce  json
// @Param endTripRequest body rest.endTripRequest true "end trip request"
// @Param api_key query string true "api_key"
// @Success 200 {object} rest.endTripResponse
// @Failure 404 {object} rest.errorRespose
// @Failure 400 {object} rest.errorRespose
// @Failure 500 {object} rest.errorRespose
// @Router /auth/user/end-trip [put]
func (api *apiDetails) endTrip(c *gin.Context) {
	req := &endTripRequest{}
	err := c.BindJSON(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = api.app.EndTrip(c, req.UserID, req.ScooterID)
	if err != nil {
		errStatusCode := getErrHTTPStatusCode(err)
		createErrorResponse(c, errStatusCode, err.Error())
		return
	}

	resp := endTripResponse{
		UserID:    req.UserID,
		ScooterID: req.ScooterID,
	}

	c.IndentedJSON(http.StatusOK, resp)
	c.Done()
}

// saveScooterTripEvent godoc
// @Summary saves the trip event generated by scooter
// @Description saves the events generated by scooter when trip is started, ended and during the trip
// @Tags scooter-api
// @Accept  json
// @Produce  json
// @Param saveScooterTripEventRequest body rest.saveScooterTripEventRequest true "save trip event request"
// @Param api_key query string true "api_key"
// @Success 200 {object} rest.saveScooterTripEventResponse
// @Failure 404 {object} rest.errorRespose
// @Failure 400 {object} rest.errorRespose
// @Failure 500 {object} rest.errorRespose
// @Router /auth/scooter/trip-event [post]
func (api *apiDetails) saveScooterTripEvent(c *gin.Context) {
	req := &saveScooterTripEventRequest{}
	err := c.BindJSON(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(req)
	if err != nil {
		createErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !domain.IsValidTripEventType(req.Type) {
		createErrorResponse(c, http.StatusBadRequest, "invalid event type")
		return
	}

	tripEvent := &domain.TripEvent{
		UserID:    req.UserID,
		ScooterID: req.ScooterID,
		Location: domain.GeoLocation{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		Type:      domain.TripEventType(req.Type),
		CreatedAt: req.CreatedAt.UTC(),
	}

	err = api.app.SaveScooterTripEvent(c, tripEvent)
	if err != nil {
		errStatusCode := getErrHTTPStatusCode(err)
		createErrorResponse(c, errStatusCode, err.Error())
		return
	}

	resp := saveScooterTripEventResponse{
		Success: true,
	}

	c.IndentedJSON(http.StatusCreated, resp)
	c.Done()
}
