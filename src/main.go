package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/myroslavve/genesis-test-case/src/api"
	"github.com/myroslavve/genesis-test-case/src/db"
	"github.com/myroslavve/genesis-test-case/src/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load the .env file
	envPath := filepath.Join(".", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database
	db.InitDB()
	api.SetDB()
	services.InitEmailService()

	// Run migrations
	runMigrations()

	// Initialize routes
	r := api.InitializeRoutes()

	// Start the email service to send exchange rates every 24 hours
	go services.SendExchangeRates(24 * time.Hour)

	// Start the server on port 8080
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}

func runMigrations() {
	dbHost := os.Getenv("MONGO_HOST")
	dbPort := os.Getenv("MONGO_PORT")
	dbUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	dbPass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	dbName := "subscriptiondb" // use your actual database name
	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", dbUser, dbPass, dbHost, dbPort, dbName)

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Ping the MongoDB server to ensure a successful connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a new MongoDB driver instance
	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: dbName})
	if err != nil {
		log.Fatalf("Failed to create MongoDB driver instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Force the migration version if the database is dirty
	if err := m.Force(1); err != nil {
		log.Fatalf("Failed to force migration version: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations ran successfully")
}
