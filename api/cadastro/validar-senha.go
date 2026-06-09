package cadastro

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/passwords"
)

func ValidarSenha(ctx context.Context, email, senha string) (bool, string) {
	var cadastro mongo.Cadastro

	err := mongo.Dns.Collection("cadastros").FindOne(ctx, bson.M{"email": email}).Decode(&cadastro)
	if err != nil {
		return false, ""
	}

	isValid := passwords.VerifyBCrypt(senha, cadastro.Senha)
	log.Printf("Validando senha para o email %q, senha: %q, válida: %v", email, senha, isValid)
	log.Println(cadastro)

	return isValid, cadastro.Email
}
