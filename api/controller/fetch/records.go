package fetch

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/records"
	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

func Records(c *gin.Context) {
	ctx := c.Request.Context()

	zona := c.Query("zone")
	if zona == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Zona é obrigatória"})

		return
	}

	nivel, err := users.FetchNivel(ctx, zona, c.GetString("email"), c.GetString("level"))
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar permissão"})

		return
	}

	if nivel == "" {
		c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para visualizar registros dessa zona"})

		return
	}

	lista, soa, err := records.Fetch(ctx, util.PdnsClient(), zona)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar registros no servidor DNS"})

		return
	}

	c.JSON(200, gin.H{"record": lista, "soa": soa, "nivel": nivel})
}
