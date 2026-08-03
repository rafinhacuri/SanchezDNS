package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Insert(ctx context.Context, zona, permissao, email string) error {
	now := time.Now()

	outra := "escrita"
	if permissao == "escrita" {
		outra = "leitura"
	}

	update := bson.M{
		"$addToSet": bson.M{permissao: email},
		"$set":      bson.M{"updatedAt": now},
		"$setOnInsert": bson.M{
			"createdAt": now,
			outra:       []string{},
		},
	}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := mongo.Dns.Collection("users").UpdateOne(ctx, bson.M{"zona": zona}, update, opts)

	return err
}
