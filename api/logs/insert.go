package logs

import (
	"context"
	"time"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Insert(zone, username, action, details string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := &Log{
		Zone:      zone,
		Username:  username,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	_, _ = mongo.Dns.Collection("logs").InsertOne(ctx, log)
}
