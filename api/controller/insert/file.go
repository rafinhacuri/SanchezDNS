package insert

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/s3"
)

func File(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Print(err)

		c.AbortWithStatusJSON(400, gin.H{"message": "Arquivo é obrigatório"})

		return
	}

	ctx := c.Request.Context()

	key, err := s3.Save(ctx, file, "")
	if err != nil {
		log.Print(err)

		c.AbortWithStatusJSON(502, gin.H{"message": "erro ao salvar arquivo"})

		return
	}

	c.JSON(200, gin.H{"message": key})
}
