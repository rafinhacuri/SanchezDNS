package mongo

import (
	"context"
	"time"
)

type Log struct {
	Usuario string    `bson:"usuario" json:"usuario"`
	Acao    string    `bson:"acao"    json:"acao"`
	Ip      string    `bson:"ip"      json:"ip"`
	Data    time.Time `bson:"data"    json:"data"`
}

func InsertLog(username, action, ip string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := &Log{
		Usuario: username,
		Acao:    action,
		Ip:      ip,
		Data:    time.Now(),
	}

	_, _ = Dns.Collection("logs").InsertOne(ctx, log)
}
