package insert

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/s3"
)

func File(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "api.invalid_request"})

		return
	}

	ctx := c.Request.Context()

	key, err := s3.Save(ctx, file, "")
	if err != nil {
		c.AbortWithStatusJSON(502, gin.H{"message": err.Error()})

		return
	}

	c.JSON(200, gin.H{"message": key})
}
