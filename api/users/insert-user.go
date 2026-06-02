//nolint:contextcheck
package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func InsertUser(ctx context.Context, zona, permissao, idcbpf, user string) (string, error) {
	collection := mongo.Dns.Collection("users")

	filter := bson.M{"zona": zona}

	var existing struct {
		Zona    string   `bson:"zona"`
		Leitura []string `bson:"leitura"`
		Escrita []string `bson:"escrita"`
	}

	err := collection.FindOne(ctx, filter).Decode(&existing)

	roleField := func() string {
		if permissao == "leitura" {
			return "leitura"
		}

		return "escrita"
	}()

	if err == nil {
		if slices.Contains(existing.Leitura, idcbpf) || slices.Contains(existing.Escrita, idcbpf) {
			return "", errors.New("usuário já cadastrado na zona com alguma permissão")
		}

		update := bson.M{
			"$addToSet": bson.M{roleField: idcbpf},
		}

		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err.Error())

			return "", errors.New("falha ao atualizar usuário")
		}

		_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"updatedAt": time.Now()}})
		if err != nil {
			log.Println(err.Error())

			return "", errors.New("falha ao atualizar usuário")
		}

		return "user updated", nil
	}

	newDoc := bson.M{
		"zona":      zona,
		"leitura":   []string{},
		"escrita":   []string{},
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	if permissao == "leitura" {
		newDoc["leitura"] = []string{idcbpf}
	} else {
		newDoc["escrita"] = []string{idcbpf}
	}

	_, err = collection.InsertOne(ctx, newDoc)
	if err != nil {
		log.Println(err.Error())

		return "", errors.New("falha ao inserir usuário")
	}

	go logs.InsertLog(
		zona,
		user,
		"insert_user",
		fmt.Sprintf("Inserido usuário %s na zona %s com permissão %s", idcbpf, zona, permissao))

	return "user inserted", nil
}
