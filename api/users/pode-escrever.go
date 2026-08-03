package users

import "context"

func PodeEscrever(ctx context.Context, zona, email, level string) (bool, error) {
	if level == "admin" {
		return true, nil
	}

	permissao, err := FetchPermissao(ctx, zona, email)
	if err != nil {
		return false, err
	}

	return permissao == "escrita", nil
}
