//nolint:dupl
package mongo

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Solicitacao struct {
	Id        bson.ObjectID `bson:"_id"       json:"id"`
	Nome      string        `bson:"nome"      json:"nome"`
	Email     string        `bson:"email"     json:"email"`
	Senha     string        `bson:"senha"     json:"senha"`
	Foto      string        `bson:"foto"      json:"foto"`
	Status    string        `bson:"status"    json:"status"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

func (s *Solicitacao) Validate() error {
	if s.Nome == "" {
		return errors.New("schema.nome_obrigatorio")
	}

	if s.Email == "" {
		return errors.New("schema.email_obrigatorio")
	}

	if s.Senha == "" {
		return errors.New("schema.senha_obrigatoria")
	}

	if s.Foto == "" {
		return errors.New("schema.foto_obrigatoria")
	}

	if s.Status == "" {
		s.Status = "pendente"
	}

	if s.Id.IsZero() {
		s.Id = bson.NewObjectID()
	}

	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}

	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = time.Now()
	}

	return nil
}

func setupSolicitacao() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	solicitacoesColl := Dns.Collection("solicitacoes")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetName("email_unique").SetUnique(true),
		},
	}

	_, err := solicitacoesColl.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Println("Error creating solicitacao indexes:", err.Error())
	}
}
