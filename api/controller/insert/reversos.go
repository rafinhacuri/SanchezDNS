package insert

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
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

	var body struct {
		Reversos []string `binding:"required,min=1" json:"reversos"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Selecione ao menos um reverso"})

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

	httpc := util.PdnsClient()

	ausentes, err := records.FetchReversosAusentes(ctx, httpc, zona)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar os reversos no servidor DNS"})

		return
	}

	selecionados := records.FiltrarReversos(ausentes, body.Reversos)

	if len(selecionados) == 0 {
		c.JSON(200, gin.H{"message": "Nenhum reverso ausente para criar"})

		return
	}

	err = records.InsertReversos(ctx, httpc, selecionados)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao criar os reversos no servidor DNS"})

		return
	}

	go logs.InsertLog(
		zona,
		email,
		"insert_reverses",
		fmt.Sprintf("Criados %d registros reversos ausentes na zona %s", len(selecionados), zona))

	c.JSON(200, gin.H{"message": fmt.Sprintf("%d reverso(s) criado(s) com sucesso!", len(selecionados))})
}
