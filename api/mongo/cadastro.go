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

type Cadastro struct {
	Id        bson.ObjectID `bson:"_id"       json:"id"`
	Nome      string        `bson:"nome"      json:"nome"`
	Email     string        `bson:"email"     json:"email"`
	Senha     string        `bson:"senha"     json:"senha"`
	Foto      string        `bson:"foto"      json:"foto"`
	Level     string        `bson:"level"     json:"level"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

func (c *Cadastro) Validate() error {
	if c.Nome == "" {
		return errors.New("schema.nome_obrigatorio")
	}

	if c.Email == "" {
		return errors.New("schema.email_obrigatorio")
	}

	if c.Senha == "" {
		return errors.New("schema.senha_obrigatoria")
	}

	if c.Foto == "" {
		return errors.New("schema.foto_obrigatoria")
	}

	if c.Level == "" {
		c.Level = "member"
	}

	if c.Id.IsZero() {
		c.Id = bson.NewObjectID()
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}

	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}

	return nil
}

func setupCadastro() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cadastrosColl := Dns.Collection("cadastros")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetName("email_unique").SetUnique(true),
		},
	}

	_, err := cadastrosColl.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Println("Error creating cadastro indexes:", err.Error())
	}
}
