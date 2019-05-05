package identity

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// API is a holder for 'identity' endpoint Gin handler functions.
type API struct {
	// Provides access to the persistent 'identity' data.
	db DB
}

// NewAPI creates a default API that uses MongoDB as it's persistent store.
func NewAPI(client *mongo.Client) *API {
	db := NewMongoDB(client)
	return &API{db}
}
