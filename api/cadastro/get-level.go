package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
)

func GetLevel(ctx context.Context, email string) (string, error) {
	level, err := redis.GetLevel(ctx, email)
	if err != nil {
		var cadastro mongo.Cadastro

		err := mongo.Dns.Collection("cadastros").FindOne(ctx, bson.M{"email": email}).Decode(&cadastro)
		if err != nil {
			return "", err
		}

		level = cadastro.Level
	}

	//nolint:contextcheck
	go redis.SaveUpdateLevel(email, level)

	return level, nil
}
