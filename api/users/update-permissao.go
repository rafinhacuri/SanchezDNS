package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func UpdatePermissao(ctx context.Context, zona, origem, destino, atual, novo string) error {
	update := bson.M{
		"$pull":     bson.M{origem: atual},
		"$addToSet": bson.M{destino: novo},
		"$set":      bson.M{"updatedAt": time.Now()},
	}

	_, err := mongo.Dns.Collection("users").UpdateOne(ctx, bson.M{"zona": zona}, update)

	return err
}
