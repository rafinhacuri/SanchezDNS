package users

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Fetch(ctx context.Context, zona string) (Zone, error) {
	var zone Zone

	err := mongo.Dns.Collection("users").FindOne(ctx, bson.M{"zona": zona}).Decode(&zone)
	if err != nil {
		return Zone{}, err
	}

	return zone, nil
}
