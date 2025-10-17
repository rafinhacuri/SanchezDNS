package controllers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatisticsResponse struct {
	Zones         int    `json:"zones"`
	Records       int    `json:"records"`
	Users         int    `json:"users"`
	Uptime        string `json:"uptime"`
	Status        string `json:"status"`
	FailedQueries int    `json:"failedQueries"`
	UDPQueries    int    `json:"udpQueries"`
	TCPQueries    int    `json:"tcpQueries"`
	ServerID      string `json:"serverId"`
	StartedAt     string `json:"startedAt"`
}

func GetStatistics(ctx *gin.Context) {
	id := ctx.Query("connection")
	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	if id == "" {
		ctx.JSON(400, gin.H{"message": "connection ID is required"})
		return
	}

	connectionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid connection ID"})
		return
	}

	var connection models.Connection

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 8*time.Second)
	defer cancel()

	err = db.Database.Collection("connections").FindOne(ctxReq, bson.M{"_id": connectionID}).Decode(&connection)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "connection not found"})
		return
	}

	if !isAdmin && !slices.Contains(connection.Users, username) {
		ctx.JSON(403, gin.H{"message": "forbidden"})
		return
	}

	plainKey, err := utils.Decrypt(connection.ApiKey)
	if err != nil {
		ctx.JSON(500, gin.H{"message": fmt.Sprintf("failed to decrypt api key: %v", err)})
		return
	}

	base := utils.NormalizeBase(connection.Host)

	serverID := connection.ServerId
	if serverID == "" {
		serverID = "localhost"
	}

	httpc := resty.New().SetBaseURL(base).SetHeader("X-API-Key", plainKey).SetHeader("Accept", "application/json").SetTimeout(6 * time.Second).SetRetryCount(2)

	var statsRaw []models.PdnsStat
	statResp, err := httpc.R().SetContext(ctxReq).SetResult(&statsRaw).Get(fmt.Sprintf("/api/v1/servers/%s/statistics", serverID))

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": fmt.Sprintf("failed to reach PowerDNS: %v", err)})
		return
	}
	if statResp.IsError() {
		ctx.JSON(statResp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS statistics error: %s", statResp.Status())})
		return
	}

	statMap := make(map[string]any, len(statsRaw))
	for _, s := range statsRaw {
		statMap[s.Name] = s.Value
	}

	getInt := func(keys ...string) int {
		for _, k := range keys {
			if v, ok := statMap[k]; ok {
				switch t := v.(type) {
				case float64:
					return int(t)
				case int:
					return t
				case int64:
					return int(t)
				case string:
					if n, perr := strconv.Atoi(t); perr == nil {
						return n
					}
				}
			}
		}
		return 0
	}

	var zones []models.PdnsZone
	zonesResp, err := httpc.R().SetContext(ctxReq).SetResult(&zones).Get(fmt.Sprintf("/api/v1/servers/%s/zones", serverID))

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": fmt.Sprintf("failed to fetch zones: %v", err)})
		return
	}
	if zonesResp.IsError() {
		ctx.JSON(zonesResp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS zones error: %s", zonesResp.Status())})
		return
	}

	records := 0
	for _, z := range zones {
		var zd models.PdnsZoneDetails
		zr, zerr := httpc.R().SetContext(ctxReq).SetResult(&zd).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", serverID, z.ID))
		if zerr != nil || zr.IsError() {
			continue
		}
		for _, rr := range zd.RRsets {
			records += len(rr.Records)
		}
	}

	uptimeSec := getInt("uptime")
	resp := StatisticsResponse{
		Zones:         len(zones),
		Records:       records,
		Users:         len(connection.Users),
		Uptime:        humanUptime(uptimeSec),
		Status:        "online",
		UDPQueries:    getInt("udp-queries"),
		TCPQueries:    getInt("tcp-queries"),
		FailedQueries: getInt("servfail-answers", "nxdomain-answers", "recursion-failures"),
		ServerID:      serverID,
		StartedAt:     startedAtFromNow(uptimeSec).Format(time.RFC3339),
	}

	ctx.JSON(200, resp)
}

func humanUptime(sec int) string {
	if sec <= 0 {
		return "0s"
	}
	d := time.Duration(sec) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

func startedAtFromNow(uptimeSec int) time.Time {
	if uptimeSec <= 0 {
		return time.Now().UTC()
	}
	return time.Now().UTC().Add(-time.Duration(uptimeSec) * time.Second)
}
