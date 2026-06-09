package remove

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Zone(c *gin.Context) {
	zoneID := c.Query("id")
	if zoneID == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "ID da zona é obrigatório"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")

	_, err := zonas.DeleteZone(ctx, zoneID, email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(502, gin.H{"message": "falha ao deletar zona"})

		return
	}

	c.JSON(200, gin.H{"message": "zona excluída com sucesso"})
}
