package solicitacoes

import (
	"context"
	"sort"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Fetch(ctx context.Context, limit int, search string, skip int64) ([]mongo.Solicitacao, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter = bson.M{"$or": []bson.M{
			{"nome": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}}
	}

	total, err := mongo.Dns.Collection("solicitacoes").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetProjection(bson.M{"senha": 0}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := mongo.Dns.Collection("solicitacoes").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var solicitacoes []mongo.Solicitacao

	err = cursor.All(ctx, &solicitacoes)
	if err != nil {
		return nil, 0, err
	}

	sort.SliceStable(solicitacoes, func(i, j int) bool {
		if solicitacoes[i].Status == "pendente" {
			return true
		}

		if solicitacoes[i].Status == "rejeitada" && solicitacoes[j].Status != "pendente" {
			return true
		}

		return false
	})

	return solicitacoes, total, nil
}
