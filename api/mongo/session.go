package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Session struct {
	ID           bson.ObjectID `bson:"_id"          json:"_id"`
	SID          string        `bson:"sid"          json:"sid"`
	Active       bool          `bson:"active"       json:"active"`
	Email        string        `bson:"email"        json:"email"`
	CreatedAt    time.Time     `bson:"createdAt"    json:"createdAt"`
	LastSeen     time.Time     `bson:"lastSeen"     json:"lastSeen"`
	Ip           string        `bson:"ip"           json:"ip"`
	Os           string        `bson:"os"           json:"os"`
	Browser      string        `bson:"browser"      json:"browser"`
	Location     string        `bson:"location"     json:"location"`
	LogoutReason string        `bson:"logoutReason" json:"logoutReason"`
}

func setupSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionsColl := Dns.Collection("sessions")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "sid", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("sidUnique"),
	}

	_, err := sessionsColl.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Erro ao criar índice para sessions: %v", err)
	}
}
