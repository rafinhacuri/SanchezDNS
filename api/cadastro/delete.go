package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Delete(ctx context.Context, id bson.ObjectID) (string, error) {
	var cadastro mongo.Cadastro

	err := mongo.Dns.Collection("cadastros").FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&cadastro)
	if err != nil {
		return "", err
	}

	return cadastro.Email, nil
}
