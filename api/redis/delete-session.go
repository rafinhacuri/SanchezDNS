package redis

import (
	"context"
)

func DeleteSession(ctx context.Context, sid string) error {
	key := "sanchezdns:session:" + sid

	_, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
