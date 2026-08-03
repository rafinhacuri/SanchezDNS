package remove

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/util"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func Zona(c *gin.Context) {
	ctx := c.Request.Context()

	domain := c.Query("id")
	if domain == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "ID da zona é obrigatório"})

		return
	}

	err := zonas.Delete(ctx, util.PdnsClient(), domain)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao excluir a zona no servidor DNS"})

		return
	}

	go logs.InsertLog(domain, c.GetString("email"), "delete_zone", "Excluída a zona "+domain)

	c.JSON(200, gin.H{"message": "Zona excluída com sucesso!"})
}
