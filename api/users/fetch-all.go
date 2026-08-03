package users

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func FetchAll(ctx context.Context) (map[string]Zone, error) {
	cursor, err := mongo.Dns.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var zones []Zone

	err = cursor.All(ctx, &zones)
	if err != nil {
		return nil, err
	}

	porZona := make(map[string]Zone, len(zones))
	for _, zone := range zones {
		porZona[strings.ToLower(strings.TrimSpace(zone.Zona))] = zone
	}

	return porZona, nil
}
