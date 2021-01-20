package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetDatabase(connection string, databaseName string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(connection)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err_ping := client.Ping(context.Background(), readpref.Primary())
	if err_ping != nil {
		log.Fatal("Couldn't connect to the database", err_ping)
	} else {
		log.Println("Connected!")
	}

	return client.Database(databaseName)
}
