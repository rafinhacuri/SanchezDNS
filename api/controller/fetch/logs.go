package fetch

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
)

func Logs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search, _ := c.GetQuery("filter")

	skip := (page - 1) * limit

	ctx := c.Request.Context()

	logs, total, err := logs.FetchLogs(ctx, page, limit, search, int64(skip))
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar logs"})

		return
	}

	c.JSON(200, gin.H{"logs": logs, "total": total})
}
