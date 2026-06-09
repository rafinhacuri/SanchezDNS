package update

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/solicitacoes"
)

func SolicitacaoStatus(c *gin.Context) {
	var body struct {
		Id     bson.ObjectID `binding:"required" json:"id"`
		Status string        `binding:"required" json:"status"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Requisição inválida",
		})

		return
	}

	ctx := c.Request.Context()

	exist, err := cadastro.ExistId(ctx, body.Id)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{
			"message": "Erro ao verificar existência do cadastro",
		})

		return
	}

	if exist {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Cadastro já existe",
		})

		return
	}

	exist, err = solicitacoes.ExistId(ctx, body.Id)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{
			"message": "Erro ao verificar existência do solicitação",
		})

		return
	}

	if !exist {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Solicitação não encontrada",
		})

		return
	}

	res, err := solicitacoes.Status(ctx, body.Id, body.Status)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{
			"message": "Erro ao atualizar status da solicitação",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": res,
	})
}
