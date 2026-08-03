package cadastro

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Fetch(ctx context.Context, limit int, search string, skip int64) ([]mongo.Cadastro, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter = bson.M{"$or": []bson.M{
			{"nome": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}}
	}

	total, err := mongo.Dns.Collection("cadastros").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "level", Value: 1},
			{Key: "nome", Value: 1},
			{Key: "createdAt", Value: -1},
		}).
		SetProjection(bson.M{"senha": 0}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := mongo.Dns.Collection("cadastros").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var cadastros []mongo.Cadastro

	err = cursor.All(ctx, &cadastros)
	if err != nil {
		return nil, 0, err
	}

	return cadastros, total, nil
}
