package users

import "context"

func FetchNivel(ctx context.Context, zona, email, level string) (string, error) {
	if level == "admin" {
		return "ADMINISTRADOR", nil
	}

	permissao, err := FetchPermissao(ctx, zona, email)
	if err != nil {
		return "", err
	}

	return NivelPermissao(permissao), nil
}
