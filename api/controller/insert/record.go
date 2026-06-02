package insert

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/records"
)

type AddRecordRequest struct {
	Zone        string `json:"zone"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	VL          string `json:"vl,omitempty"`
	TTL         int    `json:"ttl"`
	Comment     string `json:"comment,omitempty"`
	SvcPriority *int   `json:"svcPriority,omitempty"`
	TargetName  string `json:"targetName,omitempty"`
	SvcParams   string `json:"svcParams,omitempty"`
	Weight      *int   `json:"weight,omitempty"`
	Port        *int   `json:"port,omitempty"`
	Target      string `json:"target,omitempty"`
	Priority    *int   `json:"priority,omitempty"`
}

func Record(c *gin.Context) {
	var request AddRecordRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": fmt.Sprintf("corpo da requisição inválido: %v", err)})

		return
	}

	if request.Comment == "" {
		request.Comment = ""
	}

	ctx := c.Request.Context()

	level := c.GetString("level")
	email := c.GetString("email")

	if level == "member" {
		coll := mongo.Dns.Collection("users")

		var perm struct {
			Escrita []string `bson:"escrita"`
		}

		err := coll.FindOne(ctx, bson.M{"zona": request.Zone}).Decode(&perm)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(403, gin.H{"message": "Zona não encontrada ou sem permissão"})

			return
		}

		userKey := strings.ToLower(strings.TrimSpace(email))
		allowed := false

		for _, u := range perm.Escrita {
			if strings.ToLower(strings.TrimSpace(u)) == userKey {
				allowed = true

				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para visualizar registros"})

			return
		}
	}

	_, err = records.InsertRecord(
		ctx,
		request.Zone,
		request.Type,
		request.Name,
		request.VL,
		request.TTL,
		request.Comment,
		request.SvcPriority,
		request.TargetName,
		request.SvcParams,
		request.Weight,
		request.Port,
		request.Target,
		request.Priority,
		email)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": fmt.Sprintf("falha ao adicionar registro: %v", err)})

		return
	}

	c.JSON(201, gin.H{"message": "registro adicionado com sucesso"})
}
