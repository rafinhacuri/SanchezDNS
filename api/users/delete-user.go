//nolint:contextcheck
package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func DeleteUser(ctx context.Context, zona, id, arrayName string, index int, user string) (string, error) {
	collection := mongo.Dns.Collection("users")

	filter := bson.M{"zona": zona}

	var existing struct {
		Zona    string   `bson:"zona"`
		Leitura []string `bson:"leitura"`
		Escrita []string `bson:"escrita"`
	}

	err := collection.FindOne(ctx, filter).Decode(&existing)
	if err != nil {
		log.Println(err.Error())

		return "", errors.New("zona não encontrada")
	}

	var idcbpf string

	leitura := existing.Leitura
	escrita := existing.Escrita

	switch arrayName {
	case "leitura":
		if index < 0 || index >= len(leitura) {
			return "", errors.New("índice fora do intervalo para leitura")
		}

		idcbpf = leitura[index]
		leitura = append(leitura[:index], leitura[index+1:]...)
	case "escrita":
		if index < 0 || index >= len(escrita) {
			return "", errors.New("índice fora do intervalo para escrita")
		}

		idcbpf = escrita[index]
		escrita = append(escrita[:index], escrita[index+1:]...)
	default:
		return "", errors.New("permissão original inválida")
	}

	update := bson.M{
		"$set": bson.M{
			"leitura":   leitura,
			"escrita":   escrita,
			"updatedAt": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err.Error())

		return "", errors.New("falha ao atualizar usuário")
	}

	go logs.InsertLog(zona, user, "delete_user", fmt.Sprintf("Removido usuário %s da zona %s", idcbpf, zona))

	return "user deleted", nil
}
