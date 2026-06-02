package controller

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/auth"
	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Logout(c *gin.Context) {
	ctx := c.Request.Context()

	sid := c.GetString("sid")

	err := auth.DeleteSession(ctx, "", sid, "user logout")
	if err != nil {
		log.Printf("Erro ao deletar sessão do usuário %s: %v", sid, err)
		c.AbortWithStatusJSON(500, gin.H{"message": "Erro interno"})

		return
	}

	var (
		siteURL string
		secure  bool
	)

	switch env.C.Production {
	case true:
		siteURL = "sanchezdns.curi.dev.br"
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

	c.JSON(200, gin.H{"message": "deslogado com sucesso"})
}
