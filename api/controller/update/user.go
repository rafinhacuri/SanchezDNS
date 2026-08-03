package update

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/users"
)

func User(c *gin.Context) {
	ctx := c.Request.Context()

	var body struct {
		ID        string `binding:"required"                       json:"id"`
		Email     string `binding:"required,email"                 json:"email"`
		Permissao string `binding:"required,oneof=escrita leitura" json:"permissao"`
		Zona      string `binding:"required"                       json:"zona"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Dados inválidos"})

		return
	}

	permissao, index, ok := users.ParseID(body.ID)
	if !ok {
		c.AbortWithStatusJSON(400, gin.H{"message": "Usuário inválido"})

		return
	}

	zone, err := users.Fetch(ctx, body.Zona)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		c.AbortWithStatusJSON(404, gin.H{"message": "Zona não encontrada"})

		return
	}

	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao buscar usuários da zona"})

		return
	}

	atual, ok := users.FetchUser(zone, permissao, index)
	if !ok {
		c.AbortWithStatusJSON(404, gin.H{"message": "Usuário não encontrado nessa zona"})

		return
	}

	if permissao == body.Permissao {
		err = users.Update(ctx, body.Zona, permissao, index, body.Email)
	} else {
		err = users.UpdatePermissao(ctx, body.Zona, permissao, body.Permissao, atual, body.Email)
	}

	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao atualizar usuário"})

		return
	}

	go logs.InsertLog(
		body.Zona,
		c.GetString("email"),
		"update_user",
		fmt.Sprintf("Atualizado usuário %s na zona %s com permissão %s", body.Email, body.Zona, body.Permissao))

	c.JSON(200, gin.H{"message": "Usuário atualizado com sucesso!"})
}
