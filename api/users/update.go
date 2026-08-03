package users

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Update(ctx context.Context, zona, permissao string, index int, email string) error {
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("%s.%d", permissao, index): email,
			"updatedAt":                            time.Now(),
		},
	}

	_, err := mongo.Dns.Collection("users").UpdateOne(ctx, bson.M{"zona": zona}, update)

	return err
}
