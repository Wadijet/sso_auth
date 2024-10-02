package database

import (
	"atk-go-server/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetInstance initializes and returns a *mongo.Client object.
// This function uses the database connection URL from the provided configuration.
//
// Parameters:
// - c: Pointer to the config.Configuration object containing configuration information.
//
// Returns:
// - *mongo.Client: The connected MongoDB client object.
//
// Notes:
// - This function will log and return an error if there is an issue during connection or connection check.
func GetInstance(c *config.Configuration) (*mongo.Client, error) {
	if c.DataBaseConnectionURL == "" {
		return nil, fmt.Errorf("database connection URL is empty")
	}

	clientOptions := options.Client().ApplyURI(c.DataBaseConnectionURL).
		SetConnectTimeout(10 * time.Second) // Set a connection timeout

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()

	err = client.Ping(ctxPing, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// CloseInstance closes the MongoDB client connection.
func CloseInstance(client *mongo.Client) error {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("Failed to disconnect MongoDB client: %v", err)
		return err
	}
	log.Println("Successfully disconnected from MongoDB")
	return nil
}
