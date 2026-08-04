package insert

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/records"
	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

func Record(c *gin.Context) {
	ctx := c.Request.Context()

	var body records.Registro

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	if body.Zone == "" || body.Type == "" || body.Name == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	if body.Type != "HTTPS" && body.Type != "SRV" && strings.TrimSpace(body.VL) == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Valor do registro é obrigatório"})

		return
	}

	email := c.GetString("email")

	pode, err := users.PodeEscrever(ctx, body.Zone, email, c.GetString("level"))
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

	err = records.Insert(ctx, httpc, body)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao adicionar o registro no servidor DNS"})

		return
	}

	err = records.InsertReverso(ctx, httpc, body.Type, records.NormalizarValor(body), body.Zone, body.Name)
	if err != nil {
		log.Println(err)
	}

	name := records.NomeCompleto(body.Zone, body.Name)

	go logs.Insert(
		name,
		email,
		"insert_record",
		fmt.Sprintf("Criado registro %s do tipo %s na zona %s", name, body.Type, body.Zone))

	c.JSON(200, gin.H{"message": "Registro adicionado com sucesso!"})
}
