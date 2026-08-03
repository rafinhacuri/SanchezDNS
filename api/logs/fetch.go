package logs

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Fetch(ctx context.Context, filter bson.M, limit, skip int64) ([]Log, error) {
	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := mongo.Dns.Collection("logs").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	list := []Log{}

	err = cursor.All(ctx, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
