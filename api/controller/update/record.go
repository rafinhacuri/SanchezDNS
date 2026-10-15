package update

import (
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

	var body struct {
		OldValue records.Registro `json:"oldValue"`
		NewValue records.Registro `json:"newValue"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	if body.OldValue.Type != body.NewValue.Type {
		c.AbortWithStatusJSON(400, gin.H{"message": "O tipo do registro não pode ser alterado"})

		return
	}

	if body.OldValue.Name != body.NewValue.Name {
		c.AbortWithStatusJSON(400, gin.H{"message": "O nome do registro não pode ser alterado"})

		return
	}

	if body.OldValue.Zone != body.NewValue.Zone {
		c.AbortWithStatusJSON(400, gin.H{"message": "A zona do registro não pode ser alterada"})

		return
	}

	if !records.TTLValido(body.NewValue.TTL) {
		c.AbortWithStatusJSON(400, gin.H{"message": "TTL precisa ser no mínimo 60 segundos"})

		return
	}

	email := c.GetString("email")

	pode, err := users.PodeEscrever(ctx, body.NewValue.Zone, email, c.GetString("level"))
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

	err = records.Update(ctx, httpc, body.NewValue, body.OldValue)
	if err != nil {
		log.Println(err)

		if motivo, ok := records.RecusaPdns(err); ok {
			c.AbortWithStatusJSON(422, gin.H{"message": "O servidor DNS recusou o registro: " + motivo})

			return
		}

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao editar o registro no servidor DNS"})

		return
	}

	mudancas, err := records.UpdateReverso(
		ctx,
		httpc,
		body.NewValue.Type,
		records.NormalizarValor(body.NewValue),
		records.NormalizarValor(body.OldValue),
		body.NewValue.Zone,
		body.NewValue.Name)
	if err != nil {
		log.Println(err)
	}

	name := records.NomeCompleto(body.NewValue.Zone, body.NewValue.Name)

	go logs.ReversosDerivados(email, fmt.Sprintf("alteração do registro %s %s", body.NewValue.Type, name), mudancas)

	go logs.Insert(
		body.NewValue.Zone,
		email,
		"edit_record",
		fmt.Sprintf(
			"Editado registro %s do tipo %s na zona %s",
			name,
			body.NewValue.Type,
			body.NewValue.Zone))

	c.JSON(200, gin.H{"message": "Registro editado com sucesso!"})
}
