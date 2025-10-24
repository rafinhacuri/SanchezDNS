package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
)

func HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	err := db.TestMongo(ctx)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok!",
	})
}
