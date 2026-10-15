package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func Update(ctx context.Context, id bson.ObjectID, nome, foto, senha string) (string, error) {
	set := bson.M{
		"nome": nome,
		"foto": foto,
	}

	if senha != "" {
		err := passwords.BCrypt(&senha)
		if err != nil {
			return "", err
		}

		set["senha"] = senha
	}

	update := bson.M{
		"$set": set,
		"$currentDate": bson.M{
			"updatedAt": true,
		},
	}

	var cadastro mongo.Cadastro

	err := mongo.Dns.Collection("cadastros").FindOneAndUpdate(ctx, bson.M{"_id": id}, update).Decode(&cadastro)
	if err != nil {
		return "", err
	}

	return cadastro.Email, nil
}
