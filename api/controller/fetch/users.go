package fetch

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func Users(c *gin.Context) {
	ctx := c.Request.Context()

	zona := c.Query("zona")
	if zona == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Zona é obrigatória"})

		return
	}

	zone, err := users.Fetch(ctx, zona)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		c.JSON(200, []users.User{})

		return
	}

	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar usuários da zona"})

		return
	}

	c.JSON(200, users.List(zone))
}
