package fetch

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/statistics"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

func Statistics(c *gin.Context) {
	ctx := c.Request.Context()

	httpc := util.PdnsClient()

	stats, err := statistics.Fetch(ctx, httpc)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar estatísticas no servidor DNS"})

		return
	}

	zones, err := statistics.FetchZonas(ctx, httpc)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar zonas no servidor DNS"})

		return
	}

	valores := statistics.Valores(stats)
	uptime := statistics.Inteiro(valores, "uptime")

	c.JSON(200, statistics.StatisticsResponse{
		Zones:      len(zones),
		Records:    statistics.ContarRecords(ctx, httpc, zones),
		Uptime:     statistics.Uptime(uptime),
		Status:     "online",
		UDPQueries: statistics.Inteiro(valores, "udp-queries"),
		TCPQueries: statistics.Inteiro(valores, "tcp-queries"),
		ServerID:   env.C.DnsServerId,
		StartedAt:  statistics.IniciadoEm(uptime).Format(time.RFC3339),
	})
}
