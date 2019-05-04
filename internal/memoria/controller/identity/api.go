package identity

import (
	"go.mongodb.org/mongo-driver/mongo"
)


type API struct {
	DB *mongo.Client
}