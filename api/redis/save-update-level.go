package redis

import (
	"context"
	"time"
)

func SaveUpdateLevel(email, level string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = client.Set(ctx, "uniposrio-fisica:level:"+email, level, time.Hour).Err()
}
