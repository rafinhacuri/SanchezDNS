package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rafinhacuri/SanchezDNS/db"
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

	domain := strings.TrimSuffix(req.Domain, ".")
	domainWithDot := domain + "."

	soaMname := strings.TrimSuffix(req.Soa.StartOfAuthority, ".")
	soaRname := strings.TrimSuffix(req.Soa.Email, ".")

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
		"name":         domainWithDot,
		"kind":         "Native",
		"soa_edit_api": "DEFAULT",
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

	dnssecPayload := map[string]any{
		"active":  true,
		"keytype": "ksk",
	}

	respDNSSEC, err := httpc.R().SetContext(ctxReq).SetBody(dnssecPayload).Post(fmt.Sprintf("/api/v1/servers/%s/zones/%s/cryptokeys", connection.ServerId, domainWithDot))

	if err != nil {
		fmt.Printf("failed to enable DNSSEC for zone %s: %v\n", domain, err)
	}
	if respDNSSEC.IsError() {
		fmt.Printf("PowerDNS DNSSEC PATCH error: status=%d, body=%s\n", respDNSSEC.StatusCode(), respDNSSEC.String())
	}

	rrsets := []map[string]any{
		{
			"name":       domainWithDot,
			"type":       "SOA",
			"ttl":        3600,
			"changetype": "REPLACE",
			"records": []map[string]any{
				{
					"content":  fmt.Sprintf("%s. %s. 1 %d %d %d %d", soaMname, soaRname, req.Soa.Refresh, req.Soa.Retry, req.Soa.Expire, req.Soa.NegativeCacheTtl),
					"disabled": false,
				},
			},
		},
	}

	body := map[string]any{
		"rrsets": rrsets,
	}

	respPatch, err := httpc.R().SetContext(ctxReq).SetBody(body).Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", connection.ServerId, domainWithDot))
	if err != nil {
		fmt.Printf("failed to create mandatory records for zone %s: %v\n", domain, err)
	}
	if respPatch.IsError() {
		fmt.Printf("PowerDNS PATCH error: status=%d, body=%s\n", respPatch.StatusCode(), respPatch.String())
	}

	log := &models.Log{
		Username:     ctx.GetString("username"),
		IdConnection: ctx.Query("connection"),
		Action:       "create_zone",
		Details:      fmt.Sprintf("Created zone %s", domain),
		Zone:         domain,
		HostServer:   connection.Host,
		CreatedAt:    time.Now(),
	}

	_, err = db.Database.Collection("logs").InsertOne(ctx.Request.Context(), log)

	if err != nil {
		ctx.JSON(500, gin.H{"message": fmt.Sprintf("failed to log zone creation: %v", err)})
		return
	}

	ctx.JSON(201, gin.H{"message": "zone created successfully"})
}

func DeleteZone(ctx *gin.Context) {
	allowed, connection := permission(ctx)
	if !allowed {
		return
	}

	zoneID := ctx.Query("id")
	if zoneID == "" {
		ctx.JSON(400, gin.H{"message": "zone ID is required"})
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

	resp, err := httpc.R().SetContext(ctxReq).Delete(fmt.Sprintf("/api/v1/servers/%s/zones/%s", connection.ServerId, zoneID))

	if err != nil {
		ctx.JSON(502, gin.H{"message": fmt.Sprintf("failed to delete zone: %v", err)})
		return
	}

	if resp.IsError() {
		ctx.JSON(resp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS error: %v", resp.String())})
		return
	}

	log := &models.Log{
		Username:     ctx.GetString("username"),
		IdConnection: ctx.Query("connection"),
		Action:       "delete_zone",
		Details:      fmt.Sprintf("Deleted zone %s", zoneID),
		Zone:         zoneID,
		HostServer:   connection.Host,
		CreatedAt:    time.Now(),
	}

	_, err = db.Database.Collection("logs").InsertOne(ctx.Request.Context(), log)
	if err != nil {
		ctx.JSON(500, gin.H{"message": fmt.Sprintf("failed to log zone deletion: %v", err)})
		return
	}

	ctx.JSON(200, gin.H{"message": "zone deleted successfully"})
}

func UpdateSOA(ctx *gin.Context) {
	allowed, connection := permission(ctx)
	if !allowed {
		return
	}

	zoneID := ctx.Query("zone")
	if zoneID == "" {
		ctx.JSON(400, gin.H{"message": "zone ID is required"})
		return
	}

	var req models.Soa
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

	soaName := strings.TrimSuffix(req.StartOfAuthority, ".")
	soaEmail := strings.TrimSuffix(req.Email, ".")

	rrsets := []map[string]any{
		{
			"name":       zoneID,
			"type":       "SOA",
			"ttl":        3600,
			"changetype": "REPLACE",
			"records": []map[string]any{
				{
					"content":  fmt.Sprintf("%s. %s. 1 %d %d %d %d", soaName, soaEmail, req.Refresh, req.Retry, req.Expire, req.NegativeCacheTtl),
					"disabled": false,
				},
			},
		},
	}

	body := map[string]any{
		"rrsets": rrsets,
	}

	resp, err := httpc.R().SetContext(ctxReq).SetBody(body).Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", connection.ServerId, zoneID))

	if err != nil {
		ctx.JSON(502, gin.H{"message": fmt.Sprintf("failed to update SOA record: %v", err)})
		return
	}

	if resp.IsError() {
		ctx.JSON(resp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS error: %v", resp.String())})
		return
	}

	log := &models.Log{
		Username:     ctx.GetString("username"),
		IdConnection: ctx.Query("connection"),
		Action:       "update_soa",
		Details:      fmt.Sprintf("Updated SOA for zone %s", zoneID),
		Zone:         zoneID,
		HostServer:   connection.Host,
		CreatedAt:    time.Now(),
	}

	_, err = db.Database.Collection("logs").InsertOne(ctx.Request.Context(), log)

	if err != nil {
		ctx.JSON(500, gin.H{"message": fmt.Sprintf("failed to log SOA update: %v", err)})
		return
	}

	ctx.JSON(200, gin.H{"message": "SOA record updated successfully"})
}
