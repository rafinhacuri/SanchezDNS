package fetch

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
)

func Membros(c *gin.Context) {
	ctx := c.Request.Context()

	membros, err := cadastro.GetMembros(ctx)
	if err != nil {
		c.JSON(500, gin.H{"message": "Erro ao buscar membros"})
		return
	}

	c.JSON(200, membros)
}