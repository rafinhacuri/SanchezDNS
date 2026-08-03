package fetch

import (
	"log"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/util"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Zonas(c *gin.Context) {
	ctx := c.Request.Context()

	tipo := zonas.NormalizarTipo(c.Query("type"))
	if !slices.Contains([]string{"normal", "reverse", "reverse-ipv6"}, tipo) {
		c.AbortWithStatusJSON(400, gin.H{"message": "Tipo deve ser normal, reverse ou reverse-ipv6"})

		return
	}

	email := c.GetString("email")
	level := c.GetString("level")

	httpc := util.PdnsClient()

	lista, err := zonas.Fetch(ctx, httpc)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar zonas no servidor DNS"})

		return
	}

	permissoes := map[string]users.Zone{}

	if level != "admin" {
		permissoes, err = users.FetchAll(ctx)
		if err != nil {
			log.Println(err)

			c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar permissões das zonas"})

			return
		}
	}

	resultado := make([]zonas.ZoneFetch, 0, len(lista))

	for _, zona := range lista {
		nome := strings.ToLower(strings.TrimSpace(zona.Name))

		if !zonas.MatchTipo(nome, tipo) || nome == "1.1.1.in-addr.arpa." || nome == "2.2.2.in-addr.arpa." {
			continue
		}

		nivel, ok := zonas.Nivel(level, email, nome, permissoes)
		if !ok {
			continue
		}

		dnssec, err := zonas.FetchDnssec(ctx, httpc, zona)
		if err != nil {
			log.Println(err)

			c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar detalhes da zona"})

			return
		}

		resultado = append(resultado, zonas.ZoneFetch{
			Name:   zona.Name,
			Serial: zona.Serial,
			Nivel:  nivel,
			Dnssec: dnssec,
		})
	}

	c.JSON(200, zonas.ZonesResponse{Zones: resultado})
}
