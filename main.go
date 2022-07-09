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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	restApi.GracefulStopServer()
}
