package redis

import (
	"context"
)

func GetLevel(ctx context.Context, email string) (string, error) {
	key := "sanchezdns:level:" + email

	res, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
