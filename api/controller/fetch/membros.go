package fetch

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
)

func Membros(c *gin.Context) {
	ctx := c.Request.Context()

	membros, err := cadastro.GetMembros(ctx)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar membros"})

		return
	}

	c.JSON(200, membros)
}
