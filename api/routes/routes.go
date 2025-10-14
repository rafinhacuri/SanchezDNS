package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/controllers"
	"github.com/rafinhacuri/SanchezDNS/middleware"
)

func RegisterRoutes(server *gin.Engine) {
	server.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "Not Found",
		})
	})

	server.GET("/healthcheck", controllers.HealthCheck)

	server.POST("/auth", controllers.Auth)

	api := server.Group("/api", middleware.Authenticate)

	api.POST("/logout", controllers.Logout)
	api.POST("/user", controllers.InsertUser)
}
