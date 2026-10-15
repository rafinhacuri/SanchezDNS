package remove

import (
	"errors"
	"fmt"
	"log"

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

	if body.Type == "MX" && !records.PrioridadeMX(&body) {
		c.AbortWithStatusJSON(400, gin.H{"message": "Valor MX inválido: informe prioridade e destino"})

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

	err = records.Delete(ctx, httpc, body)
	if err != nil {
		log.Println(err)

		if errors.Is(err, records.ErrRegistroNaoEncontrado) {
			c.AbortWithStatusJSON(404, gin.H{"message": "Registro não encontrado"})

			return
		}

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao excluir o registro no servidor DNS"})

		return
	}

	mudancas, err := records.DeleteReverso(ctx, httpc, body.Type, body.VL, body.Zone, body.Name)
	if err != nil {
		log.Println(err)
	}

	name := records.NomeCompleto(body.Zone, body.Name)

	go logs.ReversosDerivados(email, fmt.Sprintf("exclusão do registro %s %s", body.Type, name), mudancas)

	go logs.Insert(
		body.Zone,
		email,
		"delete_record",
		fmt.Sprintf("Excluído registro %s do tipo %s na zona %s", name, body.Type, body.Zone))

	c.JSON(200, gin.H{"message": "Registro excluído com sucesso!"})
}
