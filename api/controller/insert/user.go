package insert

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

type insertUserRequest struct {
	IDCBPF    string `binding:"required"                       json:"idcbpf"`
	Permissao string `binding:"required,oneof=escrita leitura" json:"permissao"`
	Zona      string `binding:"required"                       json:"zona"`
}

func User(c *gin.Context) {
	var req insertUserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})

		return
	}

	ctx := c.Request.Context()

	idcbpf := c.GetString("idcbpf")

	_, err = users.InsertUser(ctx, req.Zona, req.Permissao, req.IDCBPF, idcbpf)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})

		return
	}

	c.JSON(200, gin.H{"message": "user inserted"})
}
