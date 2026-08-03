package users

import "context"

func PodeLer(ctx context.Context, zona, email, level string) (bool, error) {
	if level == "admin" {
		return true, nil
	}

	permissao, err := FetchPermissao(ctx, zona, email)
	if err != nil {
		return false, err
	}

	return permissao != "", nil
}
