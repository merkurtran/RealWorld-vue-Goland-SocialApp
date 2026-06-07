package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	ctx, cancle := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancle()

	// mongoURI := "mongodb://localhost:27017"
	// Good
	mongoURI := "mongodb://admin:password@mongodb"

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	if err != nil {
		fmt.Printf("error connect to db:%s\n", err.Error())
	}

	if err := Client.Ping(ctx, nil); err != nil {
		fmt.Printf("MongoDB not reachable:%s\n", err)
	}

	DB = Client.Database("social")

	fmt.Println("✅ Connected to MongoDB")
}
