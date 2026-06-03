package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func First(ctx context.Context) (bool, error) {
	count, err := mongo.Dns.Collection("cadastros").CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
