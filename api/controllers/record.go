package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/utils"
)

func GetRecords(ctx *gin.Context) {
	allowed, connection := permission(ctx)
	if !allowed {
		return
	}

	zoneID := ctx.Query("zone")
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

	resp, err := httpc.R().SetContext(ctxReq).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", connection.ServerId, zoneID))

	if err != nil {
		ctx.JSON(502, gin.H{"message": fmt.Sprintf("failed to fetch records: %v", err)})
		return
	}

	if resp.IsError() {
		ctx.JSON(resp.StatusCode(), gin.H{"message": fmt.Sprintf("PowerDNS error: %v", resp.String())})
		return
	}

	var z models.Zone
	if err := json.Unmarshal(resp.Body(), &z); err != nil {
		ctx.JSON(500, gin.H{"message": fmt.Sprintf("failed to parse PowerDNS response: %v", err)})
		return
	}

	var records []models.Simplified
	var soa *models.Soa

	for _, rr := range z.RRSets {
		if rr.Type == "SOA" && len(rr.Records) > 0 {
			parts := strings.Fields(rr.Records[0].Content)
			if len(parts) >= 7 {
				refresh, _ := strconv.Atoi(parts[3])
				retry, _ := strconv.Atoi(parts[4])
				expire, _ := strconv.Atoi(parts[5])
				negTTL, _ := strconv.Atoi(parts[6])

				soa = &models.Soa{
					StartOfAuthority: parts[0],
					Email:            parts[1],
					Refresh:          refresh,
					Retry:            retry,
					Expire:           expire,
					NegativeCacheTtl: negTTL,
				}
			}
			continue
		}

		var comments []string
		if len(rr.Comments) > 0 {
			comments = strings.Split(rr.Comments[0].Content, " | ")
		}

		for i, rec := range rr.Records {
			var priority *int
			var value string
			var comment string

			if i < len(comments) {
				comment = comments[i]
			}

			switch rr.Type {
			case "MX":
				parts := strings.Fields(rec.Content)
				if len(parts) >= 2 {
					if p, err := strconv.Atoi(parts[0]); err == nil {
						priority = &p
						value = strings.Join(parts[1:], " ")
					} else {
						value = rec.Content
					}
				} else {
					value = rec.Content
				}

			case "SRV":
				parts := strings.Fields(rec.Content)
				if len(parts) >= 4 {
					p, _ := strconv.Atoi(parts[0])
					w, _ := strconv.Atoi(parts[1])
					port, _ := strconv.Atoi(parts[2])
					t := parts[3]
					if !strings.HasSuffix(t, ".") {
						t += "."
					}
					priority = &p
					weight := w
					portVal := port
					target := t
					value = rec.Content
					records = append(records, models.Simplified{
						Zone:     z.Name,
						Type:     rr.Type,
						Name:     rr.Name,
						VL:       value,
						TTL:      rr.TTL,
						Comment:  comment,
						Priority: priority,
						Weight:   &weight,
						Port:     &portVal,
						Target:   &target,
					})
					continue
				}
				value = rec.Content

			case "HTTPS":
				parts := strings.Fields(rec.Content)
				if len(parts) >= 3 {
					svcP, _ := strconv.Atoi(parts[0])
					target := parts[1]
					params := strings.Join(parts[2:], " ")
					if !strings.HasSuffix(target, ".") {
						target += "."
					}
					svcPriority := svcP
					targetName := target
					svcParams := params
					value = rec.Content
					records = append(records, models.Simplified{
						Zone:        z.Name,
						Type:        rr.Type,
						Name:        rr.Name,
						VL:          value,
						TTL:         rr.TTL,
						Comment:     comment,
						SVCPriority: &svcPriority,
						TargetName:  &targetName,
						SVCParams:   &svcParams,
					})
					continue
				}
				value = rec.Content

			default:
				value = rec.Content
				records = append(records, models.Simplified{
					Zone:     z.Name,
					Type:     rr.Type,
					Name:     rr.Name,
					VL:       value,
					TTL:      rr.TTL,
					Comment:  comment,
					Priority: priority,
				})
			}
		}

	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})
	ctx.JSON(200, gin.H{"record": records, "soa": soa})
}

func normalizeRecordValue(req *models.AddRecordRequest) {
	if req.Type == "TXT" && req.VL != "" && !strings.HasPrefix(req.VL, "\"") {
		req.VL = fmt.Sprintf("\"%s\"", req.VL)
	}

	if (req.Type == "CNAME" || req.Type == "NS" || req.Type == "ALIAS" || req.Type == "MX") && req.VL != "" {
		if !strings.HasSuffix(req.VL, ".") {
			req.VL = req.VL + "."
		}
	}

	if req.Type == "MX" && req.Priority != nil {
		req.VL = fmt.Sprintf("%d %s", *req.Priority, req.VL)
	}

	if req.Type == "CAA" && req.VL != "" && !strings.Contains(req.VL, "issue") && !strings.Contains(req.VL, "iodef") {
		req.VL = fmt.Sprintf("0 issue \"%s\"", req.VL)
	}

	if req.Type == "SRV" {
		priority := 0
		weight := 0
		port := 0
		target := req.Target

		if req.Priority != nil {
			priority = *req.Priority
		}
		if req.Weight != nil {
			weight = *req.Weight
		}
		if req.Port != nil {
			port = *req.Port
		}

		if target != "" && !strings.HasSuffix(target, ".") {
			target += "."
		}

		req.VL = fmt.Sprintf("%d %d %d %s", priority, weight, port, target)
	}

	if req.Type == "HTTPS" {
		svcPriority := 0
		targetName := "."
		svcParams := req.SvcParams

		if req.SvcPriority != nil {
			svcPriority = *req.SvcPriority
		}
		if req.TargetName != "" {
			targetName = req.TargetName
			if !strings.HasSuffix(targetName, ".") {
				targetName += "."
			}
		}
		if svcParams == "" {
			svcParams = "alpn=\"h2\""
		}

		req.VL = fmt.Sprintf("%d %s %s", svcPriority, targetName, svcParams)
	}
}

func InsertRecord(ctx *gin.Context) {
	allowed, connection := permission(ctx)
	if !allowed {
		return
	}

	var request models.AddRecordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": fmt.Sprintf("invalid request body: %v", err)})
		return
	}

	if request.Comment == "" {
		request.Comment = "Added via SanchezDNS"
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

	if request.VL == "" {
		ctx.JSON(400, gin.H{"message": "value is required for this record type"})
		return
	}

	normalizeRecordValue(&request)

	if request.Type == "TXT" && !strings.HasPrefix(request.VL, "\"") {
		request.VL = fmt.Sprintf("\"%s\"", request.VL)
	}

	zone := strings.TrimSuffix(request.Zone, ".")
	name := request.Name

	if !strings.HasSuffix(name, ".") {
		if !strings.HasSuffix(name, zone) {
			name = fmt.Sprintf("%s.%s.", name, zone)
		} else {
			name = name + "."
		}
	}

	getResp, err := httpc.R().SetContext(ctxReq).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", connection.ServerId, request.Zone))
	if err != nil {
		ctx.JSON(502, gin.H{"message": fmt.Sprintf("failed to fetch existing records: %v", err)})
		return
	}

	var zoneData models.Zone
	if err := json.Unmarshal(getResp.Body(), &zoneData); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse PowerDNS response"})
		return
	}

	var existingRecords []struct {
		Content  string
		Disabled bool
		Comment  *string
	}

	for _, rr := range zoneData.RRSets {
		if rr.Type == request.Type && rr.Name == name {
			var comment *string
			if len(rr.Comments) > 0 {
				comment = &rr.Comments[0].Content
			}

			for _, rec := range rr.Records {
				existingRecords = append(existingRecords, struct {
					Content  string
					Disabled bool
					Comment  *string
				}{
					Content:  rec.Content,
					Disabled: rec.Disabled,
					Comment:  comment,
				})
			}
		}
	}

	existingRecords = append(existingRecords, struct {
		Content  string
		Disabled bool
		Comment  *string
	}{
		Content:  request.VL,
		Disabled: false,
		Comment: func() *string {
			if request.Comment != "" {
				c := request.Comment
				return &c
			}
			return nil
		}(),
	})

	sort.SliceStable(existingRecords, func(i, j int) bool {
		return existingRecords[i].Content < existingRecords[j].Content
	})

	var mergedRecords []map[string]any
	var allComments []string

	for _, rec := range existingRecords {
		mergedRecords = append(mergedRecords, map[string]any{
			"content":  rec.Content,
			"disabled": rec.Disabled,
		})

		if rec.Comment != nil && *rec.Comment != "" {
			allComments = append(allComments, *rec.Comment)
		} else if request.Comment != "" {
			allComments = append(allComments, request.Comment)
		} else {
			allComments = append(allComments, "Added via SanchezDNS")
		}
	}
	uniqueComments := map[string]bool{}
	var filteredComments []string
	for _, c := range allComments {
		if !uniqueComments[c] {
			uniqueComments[c] = true
			filteredComments = append(filteredComments, c)
		}
	}
	commentJoined := strings.Join(filteredComments, " | ")
	mergedComments := []map[string]any{
		{
			"content": commentJoined,
			"account": ctx.GetString("username"),
		},
	}

	resp, err := httpc.R().SetContext(ctxReq).SetBody(map[string]any{
		"rrsets": []map[string]any{
			{
				"name":       name,
				"type":       request.Type,
				"ttl":        request.TTL,
				"changetype": "REPLACE",
				"records":    mergedRecords,
				"comments":   mergedComments,
			},
		},
	}).Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", connection.ServerId, request.Zone))

	if err != nil {
		ctx.JSON(502, gin.H{"message": fmt.Sprintf("failed to insert record: %v", err)})
		return
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 201 {
		ctx.JSON(resp.StatusCode(), gin.H{"message": fmt.Sprintf("failed to insert record: %s", resp.String())})
		return
	}

	ctx.JSON(201, gin.H{"message": "record inserted successfully"})
}
