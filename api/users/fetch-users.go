package users

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

type User struct {
	IDCBPF    string `json:"idcbpf"`
	Permissao string `json:"permissao"`
	Zona      string `json:"zona"`
	ID        string `json:"id"`
}

func FetchUsers(ctx context.Context, zona string) ([]User, error) {
	collection := mongo.Dns.Collection("users")
	filter := bson.M{"zona": zona}

	var result struct {
		Zona    string   `bson:"zona"`
		Leitura []string `bson:"leitura"`
		Escrita []string `bson:"escrita"`
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Println(err.Error())

		return nil, errors.New("zona não encontrada")
	}

	var users []User
	for index, id := range result.Leitura {
		users = append(users, User{
			IDCBPF:    id,
			Permissao: "leitura",
			Zona:      result.Zona,
			ID:        fmt.Sprintf("leitura-%d", index),
		})
	}

	for index, id := range result.Escrita {
		users = append(users, User{
			IDCBPF:    id,
			Permissao: "escrita",
			Zona:      result.Zona,
			ID:        fmt.Sprintf("escrita-%d", index),
		})
	}

	return users, nil
}
