package update

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
)

func Cadastro(c *gin.Context) {
	var body struct {
		ID    bson.ObjectID `binding:"required" json:"id"`
		Nome  string        `binding:"required" json:"nome"`
		Foto  string        `binding:"required" json:"foto"`
		Senha string        `json:"senha"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Requisicao invalida"})

		return
	}

	ctx := c.Request.Context()

	exist, err := cadastro.ExistId(ctx, body.ID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Erro ao verificar existencia do cadastro"})

		return
	}

	if !exist {
		c.AbortWithStatusJSON(404, gin.H{"error": "Cadastro nao encontrado"})

		return
	}

	email, err := cadastro.Update(ctx, body.ID, body.Nome, body.Foto, body.Senha)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Erro ao atualizar cadastro"})

		return
	}

	detalhe := "Atualizado o cadastro de " + email
	if body.Senha != "" {
		detalhe += ", com troca de senha"
	}

	go logs.Insert("", c.GetString("email"), "update_account", detalhe)

	c.JSON(200, gin.H{"message": "Cadastro atualizado com sucesso"})
}
