package update

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/util"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Soa(c *gin.Context) {
	ctx := c.Request.Context()

	zoneID := c.Query("zone")
	if zoneID == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "ID da zona é obrigatório"})

		return
	}

	var body zonas.Soa

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	err = zonas.UpdateSoa(ctx, util.PdnsClient(), zoneID, body)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao atualizar o SOA no servidor DNS"})

		return
	}

	go logs.InsertLog(zoneID, c.GetString("email"), "update_soa", "Atualizado registro SOA para a zona "+zoneID)

	c.JSON(200, gin.H{"message": "Registro SOA atualizado com sucesso!"})
}
