package solicitacoes

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

	solicitacao := mongo.Solicitacao{
		Email: email,
		Senha: senha,
		Nome:  nome,
		Foto:  foto,
	}

	err = solicitacao.Validate()
	if err != nil {
		return err
	}

	_, err = mongo.Dns.Collection("solicitacoes").InsertOne(ctx, solicitacao)

	return err
}
