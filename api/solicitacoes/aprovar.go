package solicitacoes

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Aprovar(ctx context.Context, email string) (string, error) {
	var solicitacao mongo.Solicitacao

	err := mongo.Dns.Collection("solicitacoes").FindOne(ctx, bson.M{"email": email}).Decode(&solicitacao)
	if err != nil {
		return "", err
	}

	cadastro := mongo.Cadastro{
		Email: solicitacao.Email,
		Senha: solicitacao.Senha,
		Nome:  solicitacao.Nome,
		Foto:  solicitacao.Foto,
		Level: "member",
	}

	err = cadastro.Validate()
	if err != nil {
		return "", err
	}

	_, err = mongo.Dns.Collection("cadastros").InsertOne(ctx, cadastro)
	if err != nil {
		return "", err
	}

	_, err = mongo.Dns.Collection("solicitacoes").DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return "", err
	}

	return "Solicitação aprovada com sucesso", nil
}
