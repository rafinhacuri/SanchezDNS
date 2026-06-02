package fetch

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/s3"
)

func File(c *gin.Context) {
	file := c.Param("id")

	ctx := c.Request.Context()

	binaryPhoto, contentType, err := s3.Read(ctx, file)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Failed to read file"})

		return
	}

	filename := strings.TrimSuffix(file, filepath.Ext(file))

	c.Header("Content-Disposition", "inline; filename="+filename)

	c.Data(200, contentType, binaryPhoto)
}
