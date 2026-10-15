package solicitacoes

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func SenhaPendente(ctx context.Context, email, senha string) (bool, error) {
	var solicitacao mongo.Solicitacao

	err := mongo.Dns.Collection("solicitacoes").
		FindOne(ctx, bson.M{"email": email, "status": "pendente"}).
		Decode(&solicitacao)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return passwords.VerifyBCrypt(senha, solicitacao.Senha), nil
}
