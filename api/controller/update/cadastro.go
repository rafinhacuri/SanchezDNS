package update

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
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

	err = cadastro.Update(ctx, body.ID, body.Nome, body.Foto, body.Senha)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Erro ao atualizar cadastro"})

		return
	}

	c.JSON(200, gin.H{"message": "Cadastro atualizado com sucesso"})
}
