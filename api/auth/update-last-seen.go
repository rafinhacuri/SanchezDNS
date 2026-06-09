package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
)

func UpdateLastSeen(sid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionsColl := mongo.Dns.Collection("sessions")

	lastSeen := time.Now().UTC()

	sessionsColl.FindOneAndUpdate(ctx, bson.M{"sid": sid}, bson.M{"$set": bson.M{"lastSeen": lastSeen}})

	redis.SaveUpdateSession(sid)
}
