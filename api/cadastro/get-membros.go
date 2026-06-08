package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
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

	defer func() {
		_ = cursor.Close(ctx)
	}()

	for cursor.Next(ctx) {
		var membro Membro

		err := cursor.Decode(&membro)
		if err != nil {
			return nil, err
		}

		membros = append(membros, membro)
	}

	return membros, nil
}
