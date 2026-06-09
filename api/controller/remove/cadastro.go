package remove

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func Cadastro(c *gin.Context) {
	var bady struct {
		Id bson.ObjectID `binding:"required" json:"id"`
	}

	err := c.ShouldBindJSON(&bady)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	ctx := c.Request.Context()

	email, err := cadastro.Delete(ctx, bady.Id)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "erro ao deletar cadastro"})

		return
	}

	err = users.DeleteMany(ctx, email)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"message": "erro ao deletar usuários"})

		return
	}

	err = redis.DeleteLevel(ctx, email)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"message": "erro ao deletar nível de acesso"})

		return
	}

	c.JSON(200, gin.H{"message": "cadastro deletado"})
}
