package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://root:root@mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	DB *mongo.Client
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Fatal(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	return client, nil

}
