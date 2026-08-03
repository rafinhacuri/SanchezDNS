package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func GetMembros(ctx context.Context) ([]Membro, error) {
	cursor, err := mongo.Dns.Collection("cadastros").Find(ctx, bson.M{"level": "member"})
	if err != nil {
		return nil, err
	}

	var membros []Membro

	err = cursor.All(ctx, &membros)
	if err != nil {
		return nil, err
	}

	return membros, nil
}
