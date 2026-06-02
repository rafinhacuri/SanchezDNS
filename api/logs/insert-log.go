package logs

import (
	"context"
	"time"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

type Log struct {
	Zone      string `bson:"zone"      json:"zone"`
	Username  string `bson:"username"  json:"username"`
	Action    string `bson:"action"    json:"action"`
	Details   string `bson:"details"   json:"details"`
	CreatedAt string `bson:"createdAt" json:"createdAt"`
}

type LogsResponse struct {
	Logs  []Log `json:"logs"`
	Total int64 `json:"total"`
}

func InsertLog(zone, username, action, details string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := &Log{
		Username:  username,
		Action:    action,
		Details:   details,
		Zone:      zone,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	_, _ = mongo.Dns.Collection("logs").InsertOne(ctx, log)
}
