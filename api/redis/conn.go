//nolint:gochecknoglobals
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

var client *redis.Client

func Connect() {
	opt, err := redis.ParseURL(env.C.RedisUrl)
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(opt)
}

func Test(ctx context.Context) error {
	_, err := client.Ping(ctx).Result()

	return err
}
