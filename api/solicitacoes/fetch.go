package solicitacoes

import (
	"context"
	"errors"
	"log"
	"sort"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Fetch(
	ctx context.Context,
	page,
	limit int,
	search string,
	skip int64,
) ([]mongo.Solicitacao, int64, error) {
	filter := bson.M{}

	var orFilters []bson.M
	if search != "" {
		orFilters = []bson.M{
			{"nome": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if len(orFilters) > 0 {
		filter = bson.M{"$or": orFilters}
	}

	total, err := mongo.Dns.Collection("solicitacoes").CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao contar documento")
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := mongo.Dns.Collection("solicitacoes").Find(ctx, filter, opts)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao buscar solicitações")
	}

	var solicitacoes []mongo.Solicitacao

	err = cursor.All(ctx, &solicitacoes)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao analisar solicitações")
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
