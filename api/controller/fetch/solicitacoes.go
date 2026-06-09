package fetch

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/solicitacoes"
)

type SolicitacaoResponse struct {
	Solicitacoes []mongo.Solicitacao `json:"solicitacoes"`
	Total        int64               `json:"total"`
}

func Solicitacoes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search, _ := c.GetQuery("filter")

	skip := (page - 1) * limit

	ctx := c.Request.Context()

	solicitacoes, total, err := solicitacoes.Fetch(ctx, page, limit, search, int64(skip))
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar solicitações"})

		return
	}

	c.JSON(200, SolicitacaoResponse{Solicitacoes: solicitacoes, Total: total})
}
