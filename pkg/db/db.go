package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/oyaro-tech/auth/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitMongoDB() error {
	godotenv.Load()

	mongo_uri := os.Getenv("MONGODB_URI")
	if mongo_uri == "" {
		logger.Error.Fatalln("you must set \"MONGODB_URI\" environment variable")
	}

	// Set up client options
	clientOptions := options.Client().ApplyURI(mongo_uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Ping the MongoDB server
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	logger.Debug.Println("connected to MongoDB")

	Client = client
	return nil
}
