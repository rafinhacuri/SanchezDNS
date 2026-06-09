package redis

import (
	"context"
)

func DeleteLevel(ctx context.Context, email string) error {
	key := "sanchezdns:level:" + email

	_, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
