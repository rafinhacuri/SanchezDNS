package insert

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

type insertUserRequest struct {
	Email     string `binding:"required"                       json:"email"`
	Permissao string `binding:"required,oneof=escrita leitura" json:"permissao"`
	Zona      string `binding:"required"                       json:"zona"`
}

func User(c *gin.Context) {
	var req insertUserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")

	_, err = users.InsertUser(ctx, req.Zona, req.Permissao, req.Email, email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "falha ao inserir usuário"})

		return
	}

	c.JSON(200, gin.H{"message": "user inserted"})
}
