package cadastro

import (
	"context"
	"errors"
	"log"

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
) ([]mongo.Cadastro, int64, error) {
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

	total, err := mongo.Dns.Collection("cadastros").CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao contar documento")
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "level", Value: 1},
			{Key: "nome", Value: 1},
			{Key: "createdAt", Value: -1},
		}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := mongo.Dns.Collection("cadastros").Find(ctx, filter, opts)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao buscar cadastros")
	}

	var cadastros []mongo.Cadastro

	for cursor.Next(ctx) {
		var cadastro mongo.Cadastro

		err := cursor.Decode(&cadastro)
		if err != nil {
			log.Println(err.Error())

			return nil, 0, errors.New("falha ao decodificar cadastro")
		}

		cadastros = append(cadastros, cadastro)
	}

	return cadastros, total, nil
}
