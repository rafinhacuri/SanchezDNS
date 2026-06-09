package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func ExistId(ctx context.Context, id bson.ObjectID) (bool, error) {
	count, err := mongo.Dns.Collection("cadastros").CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
