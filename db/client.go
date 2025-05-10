package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"real-time-chat-app/config"
)

var (
	Client *mongo.Client
	Ctx    = context.Background()
)

// InitializeDB initializes the MongoDB connection
func InitializeDB() error {
	// Set client options
	log.Println("Connecting to MongoDB...")
	clientOptions := options.Client().ApplyURI(config.MongoDBConfig.MongoURI)
	clientOptions.SetMaxPoolSize(uint64(config.MongoDBConfig.MaxPoolSize))
	clientOptions.SetMinPoolSize(uint64(config.MongoDBConfig.MinPoolSize))
	clientOptions.SetServerSelectionTimeout(time.Second * time.Duration(config.MongoDBConfig.Timeout))

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(Ctx, clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	if err := Client.Ping(Ctx, readpref.Primary()); err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")
	return nil
}

// CloseDB closes the MongoDB connection
func CloseDB() {
	if err := Client.Disconnect(Ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}

// GetCollection returns a MongoDB collection
func GetCollection(name string) *mongo.Collection {
	return Client.Database(config.MongoDBConfig.DatabaseName).Collection(name)
}
