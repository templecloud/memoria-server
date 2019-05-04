package persistence

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo Client

func MongoClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
