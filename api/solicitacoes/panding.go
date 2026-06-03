package solicitacoes

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Panding(ctx context.Context, email string) (bool, error) {
	count, err := mongo.Dns.Collection("solicitacoes").CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
