package update

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/controller/insert"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func SOA(c *gin.Context) {
	zoneID := c.Query("zone")
	if zoneID == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "ID da zona é obrigatório"})

		return
	}

	var req insert.Soa

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	err = req.Validate()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")

	_, err = zonas.UpdateSoa(
		ctx,
		zoneID,
		req.StartOfAuthority,
		req.Email,
		req.Refresh,
		req.Retry,
		req.Expire,
		req.NegativeCacheTtl,
		email)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(502, gin.H{"message": "falha ao atualizar registro SOA"})

		return
	}

	c.JSON(200, gin.H{"message": "registro SOA atualizado com sucesso"})
}
