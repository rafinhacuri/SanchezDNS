package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
	"github.com/rafinhacuri/SanchezDNS/api/s3"
)

func HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	err := mongo.Test(ctx)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})

		return
	}

	err = redis.Test(ctx)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})

		return
	}

	err = s3.Test(ctx)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})

		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}
