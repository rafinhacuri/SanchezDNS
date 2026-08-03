package users

import (
	"context"
	"errors"

	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func FetchPermissao(ctx context.Context, zona, email string) (string, error) {
	zone, err := Fetch(ctx, zona)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return Permissao(zone, email), nil
}
