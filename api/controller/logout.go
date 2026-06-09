package controller

import (
	"log"
	"net/url"

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

	u, err := url.Parse(siteURL)
	if err == nil {
		siteURL = u.Hostname()
	}

	switch env.C.Production {
	case true:
		secure = true
	case false:
		secure = false
	}

	c.SetCookie("sanchezdns_session_id", "", -1, "/", siteURL, secure, true)

	c.JSON(200, gin.H{"message": "deslogado com sucesso"})
}
