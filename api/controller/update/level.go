package update

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
)

func Level(c *gin.Context) {
	var body struct {
		Id bson.ObjectID `binding:"required" json:"id"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Requisição inválida"})

		return
	}

	ctx := c.Request.Context()

	exist, err := cadastro.ExistId(ctx, body.Id)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"error": "Erro ao verificar existência do usuário"})

		return
	}

	if !exist {
		c.AbortWithStatusJSON(404, gin.H{"error": "Usuário não encontrado"})

		return
	}

	email, err := cadastro.UpdateLevel(ctx, body.Id)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"error": "Erro ao atualizar nível do usuário"})

		return
	}

	err = redis.DeleteLevel(ctx, email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"error": "Erro ao deletar nível do usuário no Redis"})

		return
	}

	c.JSON(200, gin.H{"message": "Nível do usuário atualizado com sucesso"})
}
