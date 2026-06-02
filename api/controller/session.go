package controller

import (
	"github.com/gin-gonic/gin"
)

type SessionRes struct {
	Idcbpf      string `json:"idcbpf"`
	Representar bool   `json:"representar"`
	Level       string `json:"level"`
}

func Session(c *gin.Context) {
	idcbpf := c.GetString("idcbpf")
	representar := c.GetBool("representar")
	level := c.GetString("level")

	c.JSON(200, SessionRes{
		Idcbpf:      idcbpf,
		Representar: representar,
		Level:       level,
	})
}
