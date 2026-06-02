package update

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/controller/insert"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/records"
)

type EditRecordRequest struct {
	OldValue insert.AddRecordRequest `json:"oldValue"`
	NewValue insert.AddRecordRequest `json:"newValue"`
}

func Record(c *gin.Context) {
	var request EditRecordRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("Requisição inválida: %v", err.Error())})

		return
	}

	if request.OldValue.Type != request.NewValue.Type {
		c.JSON(400, gin.H{"message": "tipo não pode ser alterado"})

		return
	}

	if request.OldValue.Name != request.NewValue.Name {
		c.JSON(400, gin.H{"message": "nome do registro não pode ser alterado"})

		return
	}

	ctx := c.Request.Context()

	level := c.GetString("level")
	email := c.GetString("email")

	if level == "member" {
		coll := mongo.Dns.Collection("users")

		var perm struct {
			Escrita []string `bson:"escrita"`
		}

		err := coll.FindOne(ctx, bson.M{"zona": request.OldValue.Zone}).Decode(&perm)
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

	oldRec := records.RecordChange{
		Zona:        request.OldValue.Zone,
		Tipo:        request.OldValue.Type,
		Name:        request.OldValue.Name,
		VL:          request.OldValue.VL,
		TTL:         request.OldValue.TTL,
		Comment:     request.OldValue.Comment,
		SvcPriority: request.OldValue.SvcPriority,
		TargetName:  request.OldValue.TargetName,
		SvcParams:   request.OldValue.SvcParams,
		Weight:      request.OldValue.Weight,
		Port:        request.OldValue.Port,
		Target:      request.OldValue.Target,
		Priority:    request.OldValue.Priority,
	}

	newRec := records.RecordChange{
		Zona:        request.NewValue.Zone,
		Tipo:        request.NewValue.Type,
		Name:        request.NewValue.Name,
		VL:          request.NewValue.VL,
		TTL:         request.NewValue.TTL,
		Comment:     request.NewValue.Comment,
		SvcPriority: request.NewValue.SvcPriority,
		TargetName:  request.NewValue.TargetName,
		SvcParams:   request.NewValue.SvcParams,
		Weight:      request.NewValue.Weight,
		Port:        request.NewValue.Port,
		Target:      request.NewValue.Target,
		Priority:    request.NewValue.Priority,
	}

	_, err = records.UpdateRecord(ctx, newRec, oldRec, email)
	if err != nil {
		c.JSON(500, gin.H{"message": fmt.Sprintf("falha ao editar registro: %v", err.Error())})

		return
	}

	c.JSON(200, gin.H{"message": "registro editado com sucesso"})
}
