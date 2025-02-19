package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	// Database connection
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Database disconnection
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("go_jwt").Collection(collectionName)

	return collection
}
