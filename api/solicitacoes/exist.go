package solicitacoes

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Exist(ctx context.Context, email string) (bool, error) {
	count, err := mongo.Dns.Collection("solicitacoes").CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	count, err = mongo.Dns.Collection("cadastros").CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
