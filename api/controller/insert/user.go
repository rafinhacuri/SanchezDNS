package insert

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func User(c *gin.Context) {
	ctx := c.Request.Context()

	var body struct {
		Email     string `binding:"required,email"                 json:"email"`
		Permissao string `binding:"required,oneof=escrita leitura" json:"permissao"`
		Zona      string `binding:"required"                       json:"zona"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	zone, err := users.Fetch(ctx, body.Zona)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar usuários da zona"})

		return
	}

	if users.Permissao(zone, body.Email) != "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Usuário já cadastrado nessa zona"})

		return
	}

	err = users.Insert(ctx, body.Zona, body.Permissao, body.Email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao inserir usuário"})

		return
	}

	go logs.InsertLog(
		body.Zona,
		c.GetString("email"),
		"insert_user",
		fmt.Sprintf("Inserido usuário %s na zona %s com permissão %s", body.Email, body.Zona, body.Permissao))

	c.JSON(200, gin.H{"message": "Usuário inserido com sucesso!"})
}
