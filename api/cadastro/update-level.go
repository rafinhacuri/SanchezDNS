package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func UpdateLevel(ctx context.Context, id bson.ObjectID) (string, error) {
	var cadastro mongo.Cadastro

	err := mongo.Dns.Collection("cadastros").FindOne(ctx, bson.M{"_id": id}).Decode(&cadastro)
	if err != nil {
		return "", err
	}

	var level string

	switch cadastro.Level {
	case "admin":
		level = "member"
	case "member":
		level = "admin"
	}

	_, err = mongo.Dns.Collection("cadastros").UpdateByID(ctx, id, bson.M{"$set": bson.M{"level": level}})
	if err != nil {
		return "", err
	}

	return cadastro.Email, nil
}
