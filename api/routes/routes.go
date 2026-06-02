package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/controller"
	"github.com/rafinhacuri/SanchezDNS/api/controller/fetch"
	"github.com/rafinhacuri/SanchezDNS/api/controller/insert"
	"github.com/rafinhacuri/SanchezDNS/api/controller/remove"
	"github.com/rafinhacuri/SanchezDNS/api/controller/update"
)

func RegisterRoutes(server *gin.Engine) {
	server.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "Not Found",
		})
	})

	server.GET("/healthcheck", controller.HealthCheck)

	api := server.Group("/api")

	api.POST("/login", controller.Login)
	api.POST("/logout", controller.Logout)
	api.GET("/session", controller.Session)

	api.GET("/file/:id", fetch.File)
	api.PUT("/file", insert.File)

	api.GET("/zones", fetch.Zones)
	api.GET("/records", fetch.Records)
	api.PUT("/records", insert.Record)
	api.DELETE("/records", remove.Record)
	api.PATCH("/records", update.Record)
	api.GET("/statistics", fetch.Statistics)

	api.POST("/cadastro", insert.Cadastro)

	admin := api.Group("/")

	admin.GET("/logs", fetch.Logs)
	admin.PUT("/zone", insert.Zone)
	admin.DELETE("/zone", remove.Zone)
	admin.PATCH("/soa", update.SOA)
	admin.POST("/user", insert.User)
	admin.PATCH("/user", update.User)
	admin.DELETE("/user", remove.User)
	admin.GET("/users", fetch.Users)
}
