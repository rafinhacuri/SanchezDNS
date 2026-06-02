package fetch

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/records"
)

func Records(c *gin.Context) {
	zoneID := c.Query("zone")
	if zoneID == "" {
		c.AbortWithStatusJSON(400, gin.H{"message": "Zona é obrigatória"})

		return
	}

	ctx := c.Request.Context()

	idcbpf := c.GetString("idcbpf")
	level := c.GetString("level")

	if level == "member" {
		allowed, err := memberCanViewZone(ctx, zoneID, idcbpf)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(403, gin.H{"message": "Zona não encontrada ou sem permissão"})

			return
		}

		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para visualizar registros"})

			return
		}
	}

	records, soa, err := records.FetchRecords(ctx, zoneID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"message": fmt.Sprintf("falha ao buscar registros: %v", err.Error())})

		return
	}

	c.JSON(200, gin.H{"record": records, "soa": soa})
}

func containsCI(slice []string, value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	for _, s := range slice {
		if strings.ToLower(strings.TrimSpace(s)) == v {
			return true
		}
	}

	return false
}

func memberCanViewZone(ctx context.Context, zoneID, idcbpf string) (bool, error) {
	coll := mongo.Dns.Collection("users")

	var perm struct {
		Leitura []string `bson:"leitura"`
		Escrita []string `bson:"escrita"`
	}

	err := coll.FindOne(ctx, bson.M{"zona": zoneID}).Decode(&perm)
	if err != nil {
		return false, err
	}

	if containsCI(perm.Escrita, idcbpf) {
		return true, nil
	}

	if containsCI(perm.Leitura, idcbpf) {
		return true, nil
	}

	return false, nil
}
