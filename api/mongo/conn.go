//nolint:gochecknoglobals
package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

var (
	client *mongo.Client
	Dns    *mongo.Database
)

func Connect() {
	var err error

	client, err = mongo.Connect(options.Client().ApplyURI(env.C.MongoUrl))
	if err != nil {
		panic(err)
	}

	Dns = client.Database("dns")
}

func Test(ctx context.Context) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
