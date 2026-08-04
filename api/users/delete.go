package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Delete(ctx context.Context, zona, permissao, email string) error {
	update := bson.M{
		"$pull": bson.M{permissao: email},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	_, err := mongo.Dns.Collection("users").UpdateOne(ctx, bson.M{"zona": zona}, update)

	return err
}
