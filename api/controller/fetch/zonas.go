package fetch

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Zones(c *gin.Context) {
	zoneType := strings.ToLower(strings.TrimSpace(c.Query("type")))
	if zoneType == "" {
		zoneType = "normal"
	}

	switch zoneType {
	case "normal", "reverse", "reverse-ipv6":
	default:
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"message": "tipo deve ser um dos seguintes: normal, reverse, reverse-ipv6"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")
	level := c.GetString("level")

	filtered, err := zonas.FetchZonas(ctx, zoneType, email, level)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "falha ao buscar zonas"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"zones": filtered})
}
