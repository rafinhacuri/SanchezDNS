package insert

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
)

func Cadastro(c *gin.Context) {
	var body struct {
		Nome  string `binding:"required"       json:"nome"`
		Email string `binding:"required,email" json:"email"`
		Senha string `binding:"required,min=8" json:"senha"`
		Foto  string `binding:"required,url"   json:"foto"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Requisição inválida",
		})
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
		c.JSON(409, gin.H{
			"message": "Cadastro já existe",
		})

		return
	}

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
}
