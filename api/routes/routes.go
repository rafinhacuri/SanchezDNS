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

	api := server.Group("/")
	auth := server.Group("/", middleware.ValidateSession)

	api.POST("/login", controller.Login)
	auth.POST("/logout", controller.Logout)
	auth.GET("/session", controller.Session)

	api.GET("/file/:id", fetch.File)
	api.PUT("/file", insert.File)

	auth.GET("/zones", fetch.Zonas)
	auth.GET("/records", fetch.Records)
	auth.PUT("/records", insert.Record)
	auth.DELETE("/records", remove.Record)
	auth.PATCH("/records", update.Record)
	auth.GET("/reverses", fetch.Reversos)
	auth.PUT("/reverses", insert.Reversos)
	auth.GET("/reverses/orphans", fetch.Orfaos)
	auth.DELETE("/reverses", remove.Reverso)
	auth.GET("/statistics", fetch.Statistics)

	api.POST("/cadastro", insert.Cadastro)

	admin := auth.Group("/", middleware.AdminOnly)

	admin.PUT("/solicitacoes/status", update.SolicitacaoStatus)
	admin.GET("/logs", fetch.Logs)
	admin.PUT("/zone", insert.Zona)
	admin.DELETE("/zone", remove.Zona)
	admin.PATCH("/soa", update.Soa)
	admin.POST("/user", insert.User)
	admin.PATCH("/user", update.User)
	admin.DELETE("/user", remove.User)
	admin.GET("/users", fetch.Users)
	admin.GET("/members", fetch.Membros)
	admin.GET("/solicitacoes", fetch.Solicitacoes)
	admin.GET("/cadastros", fetch.Cadastros)
	admin.DELETE("/cadastro", remove.Cadastro)
	admin.PUT("/cadastro", update.Cadastro)
	admin.PATCH("/level", update.Level)
}
