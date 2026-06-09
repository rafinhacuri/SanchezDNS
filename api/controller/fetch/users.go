package fetch

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func Users(c *gin.Context) {
	zona := c.Query("zona")
	if zona == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "zona é obrigatória"})

		return
	}

	ctx := c.Request.Context()

	users, err := users.FetchUsers(ctx, zona)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar usuários"})

		return
	}

	c.JSON(200, users)
}
