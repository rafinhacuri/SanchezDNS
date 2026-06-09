package remove

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func User(c *gin.Context) {
	var req struct {
		Zona string `binding:"required" json:"zona"`
		ID   string `binding:"required" json:"id"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	zona := req.Zona

	parts := strings.SplitN(strings.TrimSpace(req.ID), "-", 2)
	if len(parts) != 2 {
		c.AbortWithStatusJSON(400, gin.H{"message": "formato de id inválido"})

		return
	}

	arrayName := parts[0]

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "formato de id inválido"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")

	_, err = users.DeleteUser(ctx, zona, arrayName, arrayName, index, email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "falha ao excluir usuário"})

		return
	}

	c.JSON(200, gin.H{"message": "user deleted"})
}
