package identity

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DB denotes the interface to the underlying User and Credential datastores. Currently, this is 
// implemented with a MongoDB based implementation.
type DB interface {
	// User Functions
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, email string ) (*User, error)
	HasUser(ctx context.Context, email string) (bool, error)
	// Credential Functions
	CreateCredential(ctx context.Context, credential Credential) error
	GetCredential(ctx context.Context, userID string) (*Credential, error)
}

//-------------------------------------------------------------------------------------------------
// DB Implementation - Mongo

const dbName = "memoria"
const usersCol = "users"
const credsCol = "credentials"
const userEmailField = "email"
const userIDField = "userId"

// MongoDB contains references to the resources for a MongoDB based implementation of the DB interface.
type MongoDB struct {
	client *mongo.Client
	db *mongo.Database
	users *mongo.Collection
	creds *mongo.Collection
}

// NewMongoDB is the default constructor for a MongoDB instance.
func NewMongoDB(client *mongo.Client) *MongoDB {
	db := client.Database(dbName)
	users := db.Collection(usersCol)
	creds := db.Collection(credsCol)
	return &MongoDB{client, db, users, creds}
}

//-------------------------------------------------------------------------------------------------
// Public Functions - Mongo

// CreateUser inserts a new User into the datastore.
func (db *MongoDB) CreateUser(ctx context.Context, user User) error {
	_, err := db.users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// GetUser retrieves a User from the datastore.
func (db *MongoDB) GetUser(ctx context.Context, email string) (*User, error) {
	var user User
	err := db.users.FindOne(ctx, userQuery(email)).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// HasUser return if the User is in the datastore.
func (db *MongoDB) HasUser(ctx context.Context, email string) (bool, error) {
	count, cde := db.users.CountDocuments(ctx, userQuery(email))
	if cde != nil {
		return false, cde
	}

	if count == 0 {
		return false, nil
	} else if count == 1 {
		return true, nil
	} else {
		fmt.Printf("WRN: Multiple users exist with the unique Email identity: %s", email)
		return true, nil
	}
}

// CreateCredential inserts a new Credential into the datastore.
func (db *MongoDB) CreateCredential(ctx context.Context, credential Credential) error {
	_, err := db.creds.InsertOne(ctx, credential)
	if err != nil {
		return err
	}
	return nil
}

// GetCredential retrieves a Credential from the datastore.
func (db *MongoDB) GetCredential(ctx context.Context, userID string) (*Credential, error) {
	var cred Credential
	err := db.creds.FindOne(ctx, credQuery(userID)).Decode(&cred)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

//-------------------------------------------------------------------------------------------------
// Private Functions - Mongo

func userQuery(email string) bson.M {
	return bson.M{userEmailField: bson.M{"$eq": email}}
}

func credQuery(userID string) bson.M {
	return bson.M{userIDField: bson.M{"$eq": userID}}
}
