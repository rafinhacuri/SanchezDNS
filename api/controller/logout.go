package controller

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/auth"
	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func Logout(c *gin.Context) {
	ctx := c.Request.Context()

	sid := c.GetString("sid")
	idcbpf := c.GetString("email")

	err := auth.DeleteSession(ctx, "", sid, "user logout")
	if err != nil {
		log.Printf("Erro ao deletar sessão do usuário %s: %v", sid, err)
		c.AbortWithStatusJSON(500, gin.H{"message": "api.erro_generico"})

		return
	}

	var (
		siteURL string
		secure  bool
	)

	switch env.C.Production {
	case true:
		siteURL = "cbpf.br"
		secure = true
	case false:
		if strings.Contains(env.C.SiteUrl, "cbpf.dev.br") {
			siteURL = "cbpf.dev.br"
			secure = false
		} else {
			// Localhost
			siteURL = ""
			secure = false
		}
	}

	c.SetCookie("sanchezdns_session_id", "", -1, "/", siteURL, secure, true)

	ip := c.ClientIP()
	go mongo.InsertLog(idcbpf, "Se deslogou", ip)

	c.JSON(200, gin.H{"message": "api.logout_sucesso"})
}
