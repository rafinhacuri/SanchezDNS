package middleware

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/auth"
	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func ValidateSession(c *gin.Context) {
	ctx := c.Request.Context()

	sid, err := c.Cookie("sanchezdns_session_id")
	if err != nil || sid == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Missing session ID"})

		return
	}

	email, level, err := auth.ValidateSession(ctx, sid)
	if err != nil {
		log.Print(err)

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
		c.AbortWithStatusJSON(400, gin.H{"message": "Invalid session"})

		return
	}

	if email == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Invalid session"})

		return
	}

	c.Set("sid", sid)
	c.Set("email", email)
	c.Set("level", level)
	c.Next()
}
