package fetch

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/s3"
)

func File(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	foto, err := cadastro.GetFoto(ctx, id)
	if err != nil {
		foto = id
	}

	binaryPhoto, contentType, err := s3.Read(ctx, foto)
	if err != nil {
		log.Print(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar arquivo"})

		return
	}

	filename := strings.TrimSuffix(foto, filepath.Ext(foto))

	c.Header("Content-Disposition", "inline; filename="+filename)

	c.Data(200, contentType, binaryPhoto)
}
