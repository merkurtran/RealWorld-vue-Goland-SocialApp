package tests

import (
	"Server/database"
	"Server/routes"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var app *fiber.App

func TestMain(m *testing.M) {
	// setup
	setup()

	// run tests
	code := m.Run()

	// cleanup
	cleanup()

	os.Exit(code)
}

func setup() {
	// load test env vars
	if err := godotenv.Load("../env.test"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Printf("warning: failed to load env file: %v", err)
		}
	}

	// set jws if not set
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-jwt-secret-key")
	}

	// connect to test db
	connectTestDB()

	// setup fiber app
	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// setup routes
	routes.SetupAuthRoutes(app)
	//
}

func connectTestDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// use test db
	mongouri := os.Getenv("TEST_MONGO_URI")
	if mongouri == "" {
		mongouri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	database.Client = client
	database.DB = client.Database("social_test")
}

func cleanup() {
	if database.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// drop test db
		database.DB.Drop(ctx)
		database.Client.Disconnect(ctx)
	}
}

// helper for cleanup collection
func cleanupCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collections := []string{"users"}
	for _, collection := range collections {
		database.DB.Collection(collection).Drop(ctx)
	}
}
