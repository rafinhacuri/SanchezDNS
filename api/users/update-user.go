//nolint:contextcheck
package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func UpdateUser(ctx context.Context, zona, id, permissao, idcbpf, user string) (string, error) {
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

	parts := strings.SplitN(strings.TrimSpace(id), "-", 2)
	if len(parts) != 2 {
		return "", errors.New("formato de id inválido")
	}

	originalArray := parts[0]

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("formato de id inválido")
	}

	newArray := "leitura"
	if permissao == "escrita" {
		newArray = "escrita"
	}

	if originalArray == newArray {
		update := bson.M{
			"$set": bson.M{
				fmt.Sprintf("%s.%d", originalArray, index): idcbpf,
				"updatedAt": time.Now(),
			},
		}

		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err.Error())

			return "", errors.New("falha ao atualizar usuário")
		}

		return "user updated", nil
	}

	leitura := existing.Leitura
	escrita := existing.Escrita

	switch originalArray {
	case "leitura":
		if index < 0 || index >= len(leitura) {
			return "", errors.New("índice fora do intervalo para leitura")
		}

		leitura = append(leitura[:index], leitura[index+1:]...)
	case "escrita":
		if index < 0 || index >= len(escrita) {
			return "", errors.New("índice fora do intervalo para escrita")
		}

		escrita = append(escrita[:index], escrita[index+1:]...)
	default:
		return "", errors.New("permissão original inválida")
	}

	if newArray == "leitura" {
		leitura = append(leitura, idcbpf)
	} else {
		escrita = append(escrita, idcbpf)
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

	go logs.InsertLog(
		zona,
		user,
		"update_user",
		fmt.Sprintf("Atualizado usuário %s na zona %s com permissão %s", idcbpf, zona, permissao))

	return "user updated", nil
}
