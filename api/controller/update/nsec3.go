package update

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/util"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Nsec3(c *gin.Context) {
	ctx := c.Request.Context()

	zoneID := c.Query("zone")
	if zoneID == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "ID da zona é obrigatório"})

		return
	}

	httpc := util.PdnsClient()

	detalhe, err := zonas.FetchDetalhe(ctx, httpc, zoneID)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar o DNSSEC da zona"})

		return
	}

	if !detalhe.Dnssec {
		c.AbortWithStatusJSON(400, gin.H{"message": "Ative o DNSSEC da zona antes de aplicar o NSEC3"})

		return
	}

	err = zonas.InsertNsec3(ctx, httpc, zoneID)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao aplicar o NSEC3 na zona"})

		return
	}

	err = zonas.Rectify(ctx, httpc, zoneID)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao retificar a zona após aplicar o NSEC3"})

		return
	}

	go logs.Insert(zoneID, c.GetString("email"), "update_nsec3", "Aplicado NSEC3 ("+zonas.Nsec3Param+") na zona "+zoneID)

	c.JSON(200, gin.H{"message": "NSEC3 aplicado com sucesso!"})
}
