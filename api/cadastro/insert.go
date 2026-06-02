package cadastro

import (
	"context"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func Insert(ctx context.Context, email, senha, nome, foto string) (string, error) {
	var level string

	count, err := mongo.Dns.Collection("cadastros").CountDocuments(ctx, mongo.Cadastro{Email: email})
	if err != nil {
		return "", err
	}

	if count == 0 {
		level = "admin"
	} else {
		level = "member"
	}

	err = passwords.BCrypt(&senha)
	if err != nil {
		return "", err
	}

	cadastro := mongo.Cadastro{
		Email: email,
		Senha: senha,
		Nome:  nome,
		Foto:  foto,
		Level: level,
	}

	err = cadastro.Validate()
	if err != nil {
		return "", err
	}

	_, err = mongo.Dns.Collection("cadastros").InsertOne(ctx, cadastro)
	if err != nil {
		return "", err
	}

	return "Cadastro realizado com sucesso", nil
}
