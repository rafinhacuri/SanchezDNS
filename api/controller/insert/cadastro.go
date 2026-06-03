package insert

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/solicitacoes"
)

func Cadastro(c *gin.Context) {
	var body struct {
		Nome  string `binding:"required"       json:"nome"`
		Email string `binding:"required,email" json:"email"`
		Senha string `binding:"required,min=8" json:"senha"`
		Foto  string `binding:"required"       json:"foto"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println("Cadastro request received" + err.Error())

		c.JSON(400, gin.H{
			"message": "Requisição inválida",
		})
	}

	ctx := c.Request.Context()

	first, err := cadastro.First(ctx)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Erro ao verificar cadastro",
		})

		return
	}

	if first {
		res, err := cadastro.Insert(ctx, body.Email, body.Senha, body.Nome, body.Foto)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Erro ao criar cadastro",
			})

			return
		}

		c.JSON(200, gin.H{
			"message": res,
		})

		return
	}

	exist, err := solicitacoes.Exist(ctx, body.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Erro ao verificar existência da solicitação",
		})

		return
	}

	if exist {
		c.JSON(409, gin.H{
			"message": "Usuário já existe",
		})

		return
	}

	res, err := solicitacoes.Insert(ctx, body.Email, body.Senha, body.Nome, body.Foto)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Erro ao criar solicitação",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": res,
	})
}
