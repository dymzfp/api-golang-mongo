package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_HOST = "mongodb://localhost:27017"
	DB_NAME = "golang-mongo"
)

func Connect() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI(DB_HOST)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	return client.Database(DB_NAME), nil
}