package remove

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

func Record(c *gin.Context) {
	var request insert.AddRecordRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("invalid request body: %v", err.Error())})

		return
	}

	ctx := c.Request.Context()

	level := c.GetString("level")
	idcbpf := c.GetString("idcbpf")

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

		userKey := strings.ToLower(strings.TrimSpace(idcbpf))
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

	_, err = records.DeleteRecord(
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
		idcbpf)
	if err != nil {
		c.JSON(500, gin.H{"message": fmt.Sprintf("failed to delete record: %v", err.Error())})

		return
	}

	c.JSON(200, gin.H{"message": "registro excluído com sucesso"})
}
