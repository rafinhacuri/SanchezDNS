package cadastro

import (
	"context"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func Insert(ctx context.Context, email, senha, nome, foto string) error {
	err := passwords.BCrypt(&senha)
	if err != nil {
		return err
	}

	cadastro := mongo.Cadastro{
		Email: email,
		Senha: senha,
		Nome:  nome,
		Foto:  foto,
		Level: "admin",
	}

	err = cadastro.Validate()
	if err != nil {
		return err
	}

	_, err = mongo.Dns.Collection("cadastros").InsertOne(ctx, cadastro)

	return err
}
