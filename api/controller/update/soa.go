package update

import (
	"fmt"
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
		c.AbortWithStatusJSON(400, gin.H{"message": fmt.Sprintf("corpo da requisição inválido: %v", err.Error())})

		return
	}

	err = req.Validate()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": fmt.Sprintf("erro de validação: %v", err.Error())})

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
		c.AbortWithStatusJSON(502, gin.H{"message": fmt.Sprintf("falha ao atualizar registro SOA: %v", err.Error())})

		return
	}

	c.JSON(200, gin.H{"message": "registro SOA atualizado com sucesso"})
}
