package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganeshdipdumbare/scootin-aboot-journey/api/rest"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/app"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/config"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/db/mongodb"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/domain"
	"github.com/ganeshdipdumbare/scootin-aboot-journey/testclient"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Scootin Aboot Journey API
// @version 1.0
// @description A REST server to manage scooter trips and scooter events
func main() {
	// migrate reference data - product collection
	m, err := migrate.New(
		config.Get().MigrationFilesPath,
		config.Get().MongoUri+"/"+config.Get().MongoDb)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	// complete migration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database, err := mongodb.NewMongoDB(config.Get().MongoUri, config.Get().MongoDb)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Disconnect(ctx)

	scooterApp, err := app.NewApp(database)
	if err != nil {
		log.Fatal(err)
	}

	restApi, err := rest.NewApi(scooterApp, config.Get().Port, config.Get().ApiKey)
	if err != nil {
		log.Fatal(err)
	}
	restApi.StartServer()

	startTestClients()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	restApi.GracefulStopServer()
}

func startTestClients() {
	port := config.Get().Port
	apiKey := config.Get().ApiKey

	testClientRequests := []*testclient.NewTestClientReq{
		{
			Port:   port,
			UserID: "6124edb7-5099-4147-87e6-0c9b93cd1fdb",
			CurrentLocation: &domain.GeoLocation{
				Latitude:  52.54664741862859,
				Longitude: 13.351253969417021,
			},
			TravelTime: 10 * time.Second,
			RestTime:   2 * time.Second,
			ApiKey:     apiKey,
			Radius:     1,
		},
		{
			Port:   port,
			UserID: "f3b9842c-182a-418b-92fd-95d4f46414c5",
			CurrentLocation: &domain.GeoLocation{
				Latitude:  -73.961704,
				Longitude: 40.662942,
			},
			TravelTime: 12 * time.Second,
			RestTime:   3 * time.Second,
			ApiKey:     apiKey,
			Radius:     1,
		},
		{
			Port:   port,
			UserID: "4668a2f7-c498-4e49-a82e-380c1ede0685",
			CurrentLocation: &domain.GeoLocation{
				Latitude:  -73.98241999999999,
				Longitude: 40.579505,
			},
			TravelTime: 15 * time.Second,
			RestTime:   5 * time.Second,
			ApiKey:     apiKey,
			Radius:     1,
		},
	}

	for _, v := range testClientRequests {
		req := v
		go testclient.NewTestClient(req).StartJourney()
	}

}
