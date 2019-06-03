package persistence

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDefaultMongoClient create a new default connected MongoDB client. 
func NewDefaultMongoClient() *mongo.Client {
	return NewMongoClient(NewDefaultConfig())
}

// NewMongoClient creates a new connected MongoDB client. 
func NewMongoClient(config *Config) *mongo.Client {
	// Ensure Config
	if config == nil {
		config = NewDefaultConfig()
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.uri()))
	if err != nil {
		fmt.Println("Failed to create client.")
		log.Fatal(err)
	}

	if err := client.Connect(context.TODO()); err != nil {
		fmt.Println("Failed to connect to mongodb.")
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Failed to ping mongodb.")
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
