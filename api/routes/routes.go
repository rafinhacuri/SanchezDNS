package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/controller"
)

func RegisterRoutes(server *gin.Engine) {
	server.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "Not Found",
		})
	})

	server.GET("/healthcheck", controller.HealthCheck)
}
