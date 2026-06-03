package cadastro

import (
	"context"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Membro struct {
	Email string `json:"email"`
	Nome  string `json:"nome"`
}

func GetMembros(ctx context.Context) ([]Membro, error) {
	var membros []Membro

	cursor, err := mongo.Dns.Collection("cadastros").Find(ctx, bson.M{"level": "member"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var membro Membro
		if err := cursor.Decode(&membro); err != nil {
			return nil, err
		}
		membros = append(membros, membro)
	}

	return membros, nil
}