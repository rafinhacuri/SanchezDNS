package update

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/solicitacoes"
)

func AprovarSolicitacao(c *gin.Context) {
	var body struct {
		Email string `binding:"required" json:"email"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Requisição inválida",
		})

		return
	}

	ctx := c.Request.Context()

	exist, err := cadastro.Exist(ctx, body.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Erro ao verificar existência do cadastro",
		})

		return
	}

	if exist {
		c.JSON(400, gin.H{
			"message": "Cadastro já existe",
		})

		return
	}

	res, err := solicitacoes.Aprovar(ctx, body.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Erro ao aprovar solicitação",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": res,
	})
}
