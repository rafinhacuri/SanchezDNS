package remove

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/records"
	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

func Reverso(c *gin.Context) {
	ctx := c.Request.Context()

	zona := c.Query("zone")
	nome := c.Query("name")

	if zona == "" || nome == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Zona e nome do reverso são obrigatórios"})

		return
	}

	if records.IPDoReverso(nome) == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Nome de reverso inválido"})

		return
	}

	email := c.GetString("email")

	pode, err := users.PodeEscrever(ctx, zona, email, c.GetString("level"))
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar permissão"})

		return
	}

	if !pode {
		c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para alterar registros dessa zona"})

		return
	}

	err = records.DeletePTR(ctx, util.PdnsClient(), zona, nome)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao excluir o reverso no servidor DNS"})

		return
	}

	go logs.Insert(zona, email, "delete_reverse", "Excluído reverso órfão "+nome)

	c.JSON(200, gin.H{"message": "Reverso excluído com sucesso!"})
}
