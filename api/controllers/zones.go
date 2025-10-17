package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/utils"
)

func GetZones(ctx *gin.Context) {
	allowed, connection := permission(ctx)
	if !allowed {
		return
	}

	plainKey, err := utils.Decrypt(connection.ApiKey)
	if err != nil {
		ctx.JSON(500, gin.H{"message": fmt.Sprintf("failed to decrypt api key: %v", err)})
		return
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 8*time.Second)
	defer cancel()

	base := utils.NormalizeBase(connection.Host)

	httpc := resty.New().SetBaseURL(base).SetHeader("X-API-Key", plainKey).SetHeader("Accept", "application/json").SetTimeout(6 * time.Second).SetRetryCount(2)

	var zones []models.PdnsZone
	zonesResp, err := httpc.R().SetContext(ctxReq).SetResult(&zones).Get(fmt.Sprintf("/api/v1/servers/%s/zones", connection.ServerId))

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": fmt.Sprintf("failed to fetch zones: %v", err)})
		return
	}
	if zonesResp.IsError() {
		ctx.JSON(zonesResp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS zones error: %s", zonesResp.Status())})
		return
	}

	ctx.JSON(200, gin.H{"zones": zones})
}
