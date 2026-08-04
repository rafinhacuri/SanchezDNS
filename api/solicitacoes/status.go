package solicitacoes

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Status(ctx context.Context, id bson.ObjectID, status string) error {
	var solicitacao mongo.Solicitacao

	err := mongo.Dns.Collection("solicitacoes").FindOne(ctx, bson.M{"_id": id}).Decode(&solicitacao)
	if err != nil {
		return err
	}

	if status == "aprovada" {
		cadastro := mongo.Cadastro{
			Email: solicitacao.Email,
			Senha: solicitacao.Senha,
			Nome:  solicitacao.Nome,
			Foto:  solicitacao.Foto,
			Level: "member",
		}

		err = cadastro.Validate()
		if err != nil {
			return err
		}

		_, err = mongo.Dns.Collection("cadastros").InsertOne(ctx, cadastro)
		if err != nil {
			return err
		}
	}

	_, err = mongo.Dns.Collection("solicitacoes").
		UpdateByID(
			ctx,
			solicitacao.Id,
			bson.M{"$set": bson.M{"status": status}, "$currentDate": bson.M{"updatedAt": true}},
		)

	return err
}
