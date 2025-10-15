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
	server.POST("/login", controllers.Auth)
	server.POST("/api/user", controllers.InsertUser)

	api := server.Group("/api", middleware.Authenticate)
	apiAdmin := api.Group("/", middleware.AuthenticateAdmin)
	api.GET("/check-session", middleware.CheckSession)

	api.PATCH("/user/password", controllers.ChangePassword)
	api.GET("/logs", controllers.GetLogs)
	api.GET("/statistics", controllers.GetStatistics)
	api.GET("/connections", controllers.GetConnections)
	api.GET("/connection", controllers.GetConnection)
	api.PATCH("/connection", controllers.EditConnection)
	api.DELETE("/connection", controllers.DeleteConnection)

	apiAdmin.GET("/users", controllers.GetUsers)
	apiAdmin.POST("/connections", controllers.InsertConnection)
	apiAdmin.POST("/connection/user", controllers.AddUser)
	apiAdmin.DELETE("/connection/user", controllers.RemoveUser)
}
