package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
)

func ValidateSession(ctx context.Context, sid string) (string, string, error) {
	var email string

	sessionRedis, err := redis.GetSession(ctx, sid)
	if err == nil {
		if !sessionRedis.Active {
			return "", "", errors.New("session inactive")
		}

		email = sessionRedis.Email
	} else {
		var sessionMongo mongo.Session

		sessionsColl := mongo.Dns.Collection("sessions")

		err := sessionsColl.FindOne(ctx, bson.M{"sid": sid, "active": true}).Decode(&sessionMongo)
		if err != nil {
			return "", "", err
		}

		email = sessionMongo.Email
	}

	level := "Admin"

	return email, level, err
}
