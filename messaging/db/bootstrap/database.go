package bootstrap

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDatabase(uri string, dbName string) (*MongoDatabase, error) {
	// Set a timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Printf("Connected to MongoDB at %s", uri)
	return &MongoDatabase{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (m *MongoDatabase) Database() *mongo.Database {
	return m.db
}

func (m *MongoDatabase) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("MongoDB connection closed.")
	}
}
