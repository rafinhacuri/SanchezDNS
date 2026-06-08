package fetch

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/cadastro"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

type CadastroResponse struct {
	Cadastros []mongo.Cadastro `json:"cadastros"`
	Total     int64            `json:"total"`
}

func Cadastros(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search, _ := c.GetQuery("filter")

	skip := (page - 1) * limit

	ctx := c.Request.Context()

	cadastros, total, err := cadastro.Fetch(ctx, page, limit, search, int64(skip))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"message": "Erro ao buscar cadastros"})

		return
	}

	c.JSON(200, CadastroResponse{Cadastros: cadastros, Total: total})
}
