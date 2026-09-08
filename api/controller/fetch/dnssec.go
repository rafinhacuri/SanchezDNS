package fetch

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/util"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Dnssec(c *gin.Context) {
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
		c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para visualizar o DNSSEC dessa zona"})

		return
	}

	httpc := util.PdnsClient()

	detalhe, err := zonas.FetchDetalhe(ctx, httpc, zona)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar o DNSSEC da zona"})

		return
	}

	chaves, err := zonas.FetchCryptokeys(ctx, httpc, zona)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar as chaves DNSSEC da zona"})

		return
	}

	keys := make([]zonas.DnssecKey, 0, len(chaves))

	for _, chave := range chaves {
		keys = append(keys, zonas.DnssecKey{
			KeyType:   chave.KeyType,
			Algorithm: chave.Algorithm,
			Bits:      chave.Bits,
			Active:    chave.Active,
			Published: chave.Published,
			Ds:        chave.Ds,
		})
	}

	status := zonas.DsIndisponivel

	dsPai, validado, err := zonas.FetchDsPai(ctx, zona)
	if err != nil {
		log.Println(err)
	} else {
		status = zonas.StatusDs(keys, dsPai)
	}

	c.JSON(200, zonas.DnssecStatus{
		Zone:     zona,
		Dnssec:   detalhe.Dnssec,
		Nsec3:    detalhe.Nsec3,
		DsStatus: status,
		DsPai:    dsPai,
		Validado: validado,
		Keys:     keys,
	})
}
