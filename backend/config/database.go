package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func ConnectDatabase() error {
	uri := os.Getenv("MONGODB_URI")
	databaseName := os.Getenv("MONGODB_DATABASE")

	if uri == "" {
		return fmt.Errorf("MONGODB_URI is not set")
	}

	if databaseName == "" {
		return fmt.Errorf("MONGODB_DATABASE is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	client, err := mongo.Connect(
		options.Client().
			ApplyURI(uri).
			SetServerAPIOptions(serverAPI).
			SetConnectTimeout(15 * time.Second).
			SetServerSelectionTimeout(15 * time.Second),
	)

	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	Client = client
	DB = client.Database(databaseName)

	fmt.Println("MongoDB connected successfully!")

	return nil
}