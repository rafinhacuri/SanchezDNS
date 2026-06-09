package users

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func DeleteMany(ctx context.Context, email string) error {
	_, err := mongo.Dns.Collection("users").UpdateMany(
		ctx,
		bson.M{},
		bson.M{
			"$pull": bson.M{
				"leitura": email,
				"escrita": email,
			},
		},
	)

	return err
}
