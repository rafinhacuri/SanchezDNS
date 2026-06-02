package remove

import (
	"fmt"

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

	idcbpf := c.GetString("idcbpf")

	_, err := zonas.DeleteZone(ctx, zoneID, idcbpf)
	if err != nil {
		c.AbortWithStatusJSON(502, gin.H{"message": fmt.Sprintf("falha ao deletar zona: %v", err.Error())})

		return
	}

	c.JSON(200, gin.H{"message": "zona excluída com sucesso"})
}
