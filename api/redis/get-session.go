package redis

import (
	"context"
	"encoding/json"
	"errors"
)

func GetSession(ctx context.Context, sid string) (SessionJSON, error) {
	var session SessionJSON

	key := "sanchezdns:session:" + sid

	res, err := client.Do(ctx, "JSON.GET", key).Result()
	if err != nil {
		return session, err
	}

	if res == nil {
		return session, errors.New("session not found")
	}

	str, ok := res.(string)
	if !ok {
		return session, errors.New("invalid redis response type")
	}

	err = json.Unmarshal([]byte(str), &session)
	if err != nil {
		return session, err
	}

	return session, nil
}
