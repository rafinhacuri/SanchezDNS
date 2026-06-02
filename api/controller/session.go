package controller

import (
	"github.com/gin-gonic/gin"
)

type SessionRes struct {
	Email string `json:"email"`
	Level string `json:"level"`
}

func Session(c *gin.Context) {
	email := c.GetString("email")
	level := c.GetString("level")

	c.JSON(200, SessionRes{
		Email: email,
		Level: level,
	})
}
