package controller

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"

	"github.com/rafinhacuri/SanchezDNS/api/auth"
	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/solicitacoes"
	"github.com/rafinhacuri/SanchezDNS/api/util"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"senha"`
}

func Login(c *gin.Context) {
	ctx := c.Request.Context()

	// O erro é ignorado porque se o cookie não existir, cookie será uma string vazia.
	cookie, _ := c.Cookie("sanchezdns_session_id")
	if cookie != "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "ja existe uma sessão ativa"})

		return
	}

	var credentials LoginBody

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	panding, err := solicitacoes.Panding(ctx, credentials.Email)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao verificar existência da solicitação"})

		return
	}

	if panding {
		c.AbortWithStatusJSON(401, gin.H{"message": "Seu cadastro ainda está em análise"})

		return
	}

	isValidPassword, mail := cadastro.ValidarSenha(ctx, credentials.Email, credentials.Password)
	if !isValidPassword {
		c.AbortWithStatusJSON(401, gin.H{"message": "email ou senha incorretos"})

		return
	}

	uaStr := c.GetHeader("User-Agent")
	ua := useragent.Parse(uaStr)
	ip := c.ClientIP()
	os := ua.OS
	browser := ua.Name + " " + ua.Version
	location := util.GeoLocation(ip)

	sessionID, err := auth.CreateSession(ctx, mail, ip, os, browser, location)
	if err != nil {
		log.Printf("Erro ao criar sessão para o usuário %q: %v", mail, err)
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

	c.SetCookie("sanchezdns_session_id", sessionID, 34560000, "/", siteURL, secure, true)

	c.JSON(200, gin.H{
		"message": "logado com sucesso",
	})
}
