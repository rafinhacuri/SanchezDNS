package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func Update(ctx context.Context, id bson.ObjectID, nome, foto, senha string) error {
	set := bson.M{
		"nome": nome,
		"foto": foto,
	}

	if senha != "" {
		passwords.BCrypt(&senha)

		set["senha"] = senha
	}

	update := bson.M{
		"$set": set,
		"$currentDate": bson.M{
			"updatedAt": true,
		},
	}

	_, err := mongo.Dns.Collection("cadastros").UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	return nil
}
