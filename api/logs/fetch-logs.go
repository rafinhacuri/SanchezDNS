package logs

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func FetchLogs(
	ctx context.Context,
	page,
	limit int,
	search string,
	skip int64,
) ([]bson.M, int64, error) {
	filter := bson.M{}

	var orFilters []bson.M
	if search != "" {
		orFilters = []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"action": bson.M{"$regex": search, "$options": "i"}},
			{"details": bson.M{"$regex": search, "$options": "i"}},
			{"zone": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if len(orFilters) > 0 {
		filter = bson.M{"$or": orFilters}
	}

	total, err := mongo.Dns.Collection("logs").CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao contar documento")
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := mongo.Dns.Collection("logs").Find(ctx, filter, opts)
	if err != nil {
		log.Println(err.Error())

		return nil, 0, errors.New("falha ao buscar logs")
	}

	var logsMongo []bson.M

	for cursor.Next(ctx) {
		var logEntry bson.M

		err := cursor.Decode(&logEntry)
		if err != nil {
			log.Println(err.Error())

			return nil, 0, errors.New("falha ao decodificar log")
		}

		logsMongo = append(logsMongo, logEntry)
	}

	return logsMongo, total, nil
}
