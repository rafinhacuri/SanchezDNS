package fetch

import (
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
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})

		return
	}

	c.JSON(200, users)
}
