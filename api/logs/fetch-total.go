package logs

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func FetchTotal(ctx context.Context, filter bson.M) (int64, error) {
	return mongo.Dns.Collection("logs").CountDocuments(ctx, filter)
}
