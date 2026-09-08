package insert

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/util"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Zona(c *gin.Context) {
	ctx := c.Request.Context()

	var body zonas.CreateZoneRequest

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	tipo := zonas.NormalizarTipo(body.Type)
	if !slices.Contains([]string{"normal", "reverse", "reverse-ipv6"}, tipo) {
		c.AbortWithStatusJSON(400, gin.H{"message": "Tipo deve ser normal, reverse ou reverse-ipv6"})

		return
	}

	domain := strings.TrimSuffix(strings.TrimSpace(body.Domain), ".") + "."
	if !zonas.ValidoDominio(domain, tipo) {
		c.AbortWithStatusJSON(400, gin.H{"message": "Domínio não corresponde ao tipo da zona"})

		return
	}

	httpc := util.PdnsClient()

	err = zonas.Insert(ctx, httpc, domain)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao criar a zona no servidor DNS"})

		return
	}

	err = zonas.InsertDnssec(ctx, httpc, domain)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao ativar o DNSSEC da zona"})

		return
	}

	err = zonas.InsertNsec3(ctx, httpc, domain)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao aplicar o NSEC3 na zona"})

		return
	}

	err = zonas.UpdateSoa(ctx, httpc, domain, body.Soa)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao configurar o SOA da zona"})

		return
	}

	err = zonas.Rectify(ctx, httpc, domain)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao retificar a zona"})

		return
	}

	go logs.Insert(domain, c.GetString("email"), "create_zone", fmt.Sprintf("Criada zona %s do tipo %s", domain, tipo))

	c.JSON(200, gin.H{"message": "Zona criada com sucesso!"})
}
