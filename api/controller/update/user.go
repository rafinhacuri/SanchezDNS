package update

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func User(c *gin.Context) {
	var req users.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")

	_, err = users.UpdateUser(ctx, req.Zona, req.ID, req.Permissao, req.Email, email)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao atualizar usuário"})

		return
	}

	c.JSON(200, gin.H{"message": "usuário atualizado com sucesso"})
}
