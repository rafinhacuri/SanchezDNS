package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func GetFoto(ctx context.Context, email string) (string, error) {
	var cadastro mongo.Cadastro

	err := mongo.Dns.Collection("cadastros").FindOne(ctx, bson.M{"email": email}).Decode(&cadastro)
	if err != nil {
		return "", err
	}

	return cadastro.Foto, nil
}
