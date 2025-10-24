package db

import (
	"context"
	"crypto/tls"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database
var client *mongo.Client

func InitDB(ssl bool, user, url, pass, name string) error {
	opts := options.Client().ApplyURI(url)
	if user != "" || pass != "" {
		opts.SetAuth(options.Credential{Username: user, Password: pass})
	}
	if ssl {
		opts.SetTLSConfig(&tls.Config{})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		_ = mongoClient.Disconnect(context.Background())
		return err
	}

	Database = mongoClient.Database(name)
	client = mongoClient

	return nil
}

func TestMongo(ctx context.Context) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
