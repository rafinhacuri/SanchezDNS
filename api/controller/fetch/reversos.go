//nolint:dupl // endpoint distinto do fetch.Orfaos, compartilha apenas o boilerplate de validação
package fetch

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/records"
	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

func Reversos(c *gin.Context) {
	ctx := c.Request.Context()

	zona := c.Query("zone")
	if zona == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Zona é obrigatória"})

		return
	}

	pode, err := users.PodeLer(ctx, zona, c.GetString("email"), c.GetString("level"))
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar permissão"})

		return
	}

	if !pode {
		c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para visualizar registros dessa zona"})

		return
	}

	ausentes, err := records.FetchReversosAusentes(ctx, util.PdnsClient(), zona)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar os reversos no servidor DNS"})

		return
	}

	c.JSON(200, ausentes)
}
