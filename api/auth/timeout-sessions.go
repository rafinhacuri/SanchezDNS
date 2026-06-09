package auth

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"

	"github.com/rafinhacuri/SanchezDNS/api/redis"
)

func TimeoutSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionsColl := mongo.Dns.Collection("sessions")

	cutoff := time.Now().UTC().Add(-7 * 24 * time.Hour)

	filter := bson.M{
		"lastSeen": bson.M{"$lt": cutoff},
		"active":   true,
	}

	findOpts := options.Find().SetProjection(bson.M{"sid": 1})

	cursor, err := sessionsColl.Find(ctx, filter, findOpts)
	if err != nil {
		log.Printf("Erro ao buscar sessões expiradas: %v", err)

		return
	}

	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("Erro ao fechar cursor: %v", cerr)
		}
	}()

	type sidDoc struct {
		SID string `bson:"sid"`
	}

	var sids []string

	for cursor.Next(ctx) {
		var d sidDoc

		err := cursor.Decode(&d)
		if err != nil {
			log.Printf("Erro ao decodificar sessão expirada: %v", err)

			continue
		}

		if d.SID != "" {
			sids = append(sids, d.SID)
		}
	}

	err = cursor.Err()
	if err != nil {
		log.Printf("Erro no cursor de sessões expiradas: %v", err)

		return
	}

	if len(sids) == 0 {
		return
	}

	update := bson.M{
		"$set": bson.M{
			"active":       false,
			"logoutReason": "timeout",
		},
	}

	_, err = sessionsColl.UpdateMany(ctx, bson.M{"sid": bson.M{"$in": sids}}, update)
	if err != nil {
		log.Printf("Erro ao desativar sessões expiradas: %v", err)
	}

	for _, sid := range sids {
		_ = redis.DeleteSession(ctx, sid)
	}
}
