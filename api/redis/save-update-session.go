package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func SaveUpdateSession(sid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionsColl := mongo.Dns.Collection("sessions")

	var session mongo.Session

	err := sessionsColl.FindOne(ctx, bson.M{"sid": sid}).Decode(&session)
	if err != nil {
		log.Printf("Erro ao obter sessão para SID %s: %v", sid, err)

		return
	}

	if !session.Active {
		return
	}

	sessionJson := SessionJSON{
		ID:           session.ID.Hex(),
		SID:          session.SID,
		Active:       session.Active,
		Email:        session.Email,
		CreatedAt:    session.CreatedAt,
		LastSeen:     session.LastSeen,
		Ip:           session.Ip,
		Os:           session.Os,
		Browser:      session.Browser,
		Location:     session.Location,
		LogoutReason: session.LogoutReason,
	}

	sessionMarshal, err := json.Marshal(sessionJson)
	if err != nil {
		log.Printf("Erro ao serializar sessão para SID %s: %v", sid, err)

		return
	}

	key := "uniposrio-fisica:session:" + session.SID

	pipe := client.TxPipeline()

	pipe.Do(ctx, "JSON.SET", key, "$", string(sessionMarshal))

	pipe.Expire(ctx, key, time.Hour)

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Erro ao salvar sessão no Redis para SID %s: %v", sid, err)
	}
}
