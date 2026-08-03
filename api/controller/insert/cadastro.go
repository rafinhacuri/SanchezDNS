package insert

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/solicitacoes"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

func Cadastro(c *gin.Context) {
	ctx := c.Request.Context()

	var body struct {
		Nome  string `binding:"required"       json:"nome"`
		Email string `binding:"required,email" json:"email"`
		Senha string `binding:"required,min=8" json:"senha"`
		Foto  string `binding:"required"       json:"foto"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	body.Email = strings.ToLower(body.Email)

	if !util.ValidateEmail(body.Email) {
		c.AbortWithStatusJSON(400, gin.H{"message": "Email inválido"})

		return
	}

	first, err := cadastro.First(ctx)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar cadastro"})

		return
	}

	if first {
		err = cadastro.Insert(ctx, body.Email, body.Senha, body.Nome, body.Foto)
		if err != nil {
			log.Println(err)

			c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao criar cadastro"})

			return
		}

		c.JSON(200, gin.H{"message": "Cadastro realizado com sucesso"})

		return
	}

	exist, err := solicitacoes.Exist(ctx, body.Email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar existência da solicitação"})

		return
	}

	if exist {
		c.AbortWithStatusJSON(409, gin.H{"message": "Usuário já existe"})

		return
	}

	err = solicitacoes.Insert(ctx, body.Email, body.Senha, body.Nome, body.Foto)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao criar solicitação"})

		return
	}

	c.JSON(200, gin.H{"message": "Solicitação realizada com sucesso e será analisada em breve"})
}
