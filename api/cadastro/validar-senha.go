package cadastro

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func ValidarSenha(ctx context.Context, email, senha string) (bool, string, error) {
	var cadastro mongo.Cadastro

	err := mongo.Dns.Collection("cadastros").FindOne(ctx, bson.M{"email": email}).Decode(&cadastro)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, "", nil
	}

	if err != nil {
		return false, "", err
	}

	return passwords.VerifyBCrypt(senha, cadastro.Senha), cadastro.Email, nil
}
