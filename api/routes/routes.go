package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/controller"
	"github.com/rafinhacuri/SanchezDNS/api/controller/fetch"
	"github.com/rafinhacuri/SanchezDNS/api/controller/insert"
	"github.com/rafinhacuri/SanchezDNS/api/controller/remove"
	"github.com/rafinhacuri/SanchezDNS/api/controller/update"
	"github.com/rafinhacuri/SanchezDNS/api/middleware"
)

func RegisterRoutes(server *gin.Engine) {
	server.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "Not Found",
		})
	})

	server.GET("/healthcheck", controller.HealthCheck)

	api := server.Group("/api")
	auth := api.Group("/", middleware.ValidateSession)

	api.POST("/login", controller.Login)
	auth.POST("/logout", controller.Logout)
	auth.GET("/session", controller.Session)

	api.GET("/file/:id", fetch.File)
	api.PUT("/file", insert.File)

	auth.GET("/zones", fetch.Zones)
	auth.GET("/records", fetch.Records)
	auth.PUT("/records", insert.Record)
	auth.DELETE("/records", remove.Record)
	auth.PATCH("/records", update.Record)
	auth.GET("/statistics", fetch.Statistics)

	api.POST("/cadastro", insert.Cadastro)

	admin := auth.Group("/", middleware.AdminOnly)

	admin.PUT("/solicitacoes-status", update.SolicitacaoStatus)
	admin.GET("/logs", fetch.Logs)
	admin.PUT("/zone", insert.Zone)
	admin.DELETE("/zone", remove.Zone)
	admin.PATCH("/soa", update.SOA)
	admin.POST("/user", insert.User)
	admin.PATCH("/user", update.User)
	admin.DELETE("/user", remove.User)
	admin.GET("/users", fetch.Users)
	admin.GET("/members", fetch.Membros)
	admin.GET("/solicitacoes", fetch.Solicitacoes)
	admin.GET("/cadastros", fetch.Cadastros)
	admin.DELETE("/cadastro", remove.Cadastro)
}
