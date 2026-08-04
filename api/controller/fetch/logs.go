package fetch

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
)

func Logs(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	filter := logs.Filter(c.Query("filter"))

	total, err := logs.FetchTotal(ctx, filter)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar logs"})

		return
	}

	list, err := logs.Fetch(ctx, filter, int64(limit), int64((page-1)*limit))
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar logs"})

		return
	}

	c.JSON(200, logs.LogsResponse{Logs: list, Total: total})
}
