package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/redis"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func DeleteSession(ctx context.Context, id, sid, logoutReason string) error {
	sessionsColl := mongo.Dns.Collection("sessions")

	lastSeen := time.Now().UTC()

	filter := bson.M{"sid": sid}

	if id != "" {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		filter = bson.M{"_id": oid}
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetProjection(bson.M{"sid": 1})

	var session mongo.Session

	err := sessionsColl.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": bson.M{"active": false, "lastSeen": lastSeen, "logoutReason": logoutReason}},
		opts).Decode(&session)
	if err != nil {
		return err
	}

	_ = redis.DeleteSession(ctx, session.SID)

	return nil
}
