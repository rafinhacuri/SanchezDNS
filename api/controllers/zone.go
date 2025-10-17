package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/utils"
)

func CreateZone(ctx *gin.Context) {
	allowed, connection := permission(ctx)
	if !allowed {
		return
	}

	var req models.CreateZoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"message": fmt.Sprintf("invalid request body: %v", err)})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": fmt.Sprintf("validation error: %v", err)})
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

	zonePayload := map[string]any{
		"name":         req.Name,
		"kind":         req.Kind,
		"soa_edit_api": req.SoaEditApi,
	}

	if len(req.Masters) > 0 {
		zonePayload["masters"] = req.Masters
	}
	if len(req.AlsoNotify) > 0 {
		zonePayload["also_notify"] = req.AlsoNotify
	}

	resp, err := httpc.R().SetContext(ctxReq).SetBody(zonePayload).Post(fmt.Sprintf("/api/v1/servers/%s/zones", connection.ServerId))

	if err != nil {
		ctx.JSON(502, gin.H{"message": fmt.Sprintf("failed to create zone: %v", err)})
		return
	}

	if resp.IsError() {
		ctx.JSON(resp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS error: %v", resp.String())})
		return
	}

	ctx.JSON(201, gin.H{"message": "zone created successfully"})
}
