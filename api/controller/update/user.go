package update

import (
	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func User(c *gin.Context) {
	var req users.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})

		return
	}

	ctx := c.Request.Context()

	idcbpf := c.GetString("idcbpf")

	_, err = users.UpdateUser(ctx, req.Zona, req.ID, req.Permissao, req.IDCBPF, idcbpf)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})

		return
	}

	c.JSON(200, gin.H{"message": "user updated"})
}
