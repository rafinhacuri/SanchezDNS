//nolint:tagliatelle, gocognit
package fetch

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

type PdnsZoneDetails struct {
	Name   string      `json:"name"`
	RRsets []PdnsRRSet `json:"rrsets"`
}

type PdnsRRSet struct {
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	TTL     int            `json:"ttl"`
	Records []PdnsRRRecord `json:"records"`
}

type PdnsRRRecord struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

type PdnsZone struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Serial     int64  `json:"serial"`
	URL        string `json:"url"`
	SoaEditApi string `json:"soa_edit_api"`
}

type PdnsStat struct {
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type StatisticsResponse struct {
	Zones      int    `json:"zones"`
	Records    int    `json:"records"`
	Uptime     string `json:"uptime"`
	Status     string `json:"status"`
	UDPQueries int    `json:"udpQueries"`
	TCPQueries int    `json:"tcpQueries"`
	ServerID   string `json:"serverId"`
	StartedAt  string `json:"startedAt"`
}

func Statistics(c *gin.Context) {
	ctx := c.Request.Context()

	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	var statsRaw []PdnsStat

	statResp, err := httpc.R().
		SetContext(ctx).
		SetResult(&statsRaw).
		Get(fmt.Sprintf("/api/v1/servers/%s/statistics", env.C.DnsServerId))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"message": "falha ao alcançar PowerDNS"})

		return
	}

	if statResp.IsStatusFailure() {
		c.AbortWithStatusJSON(500, gin.H{"message": "Erro ao obter estatísticas"})

		return
	}

	statMap := make(map[string]string, len(statsRaw))

	for _, s := range statsRaw {
		if len(s.Value) == 0 {
			statMap[s.Name] = ""

			continue
		}

		var asString string

		err = json.Unmarshal(s.Value, &asString)
		if err == nil {
			statMap[s.Name] = asString

			continue
		}

		var asBool bool

		err = json.Unmarshal(s.Value, &asBool)
		if err == nil {
			if asBool {
				statMap[s.Name] = "true"
			} else {
				statMap[s.Name] = "false"
			}

			continue
		}

		var asInt int64

		err = json.Unmarshal(s.Value, &asInt)
		if err == nil {
			statMap[s.Name] = strconv.FormatInt(asInt, 10)

			continue
		}

		var asFloat float64

		err = json.Unmarshal(s.Value, &asFloat)
		if err == nil {
			if asFloat == float64(int64(asFloat)) {
				statMap[s.Name] = strconv.FormatInt(int64(asFloat), 10)
			} else {
				statMap[s.Name] = strconv.FormatFloat(asFloat, 'f', -1, 64)
			}

			continue
		}

		statMap[s.Name] = string(s.Value)

		continue
	}

	getInt := func(keys ...string) int {
		for _, k := range keys {
			if v, ok := statMap[k]; ok {
				n, perr := strconv.Atoi(v)
				if perr == nil {
					return n
				}

				f, ferr := strconv.ParseFloat(v, 64)
				if ferr == nil {
					return int(f)
				}
			}
		}

		return 0
	}

	var zones []PdnsZone

	zonesResp, err := httpc.R().
		SetContext(ctx).
		SetResult(&zones).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"message": "falha ao buscar zonas"})

		return
	}

	if zonesResp.IsStatusFailure() {
		c.AbortWithStatusJSON(500, gin.H{"message": "Erro nas zonas do PowerDNS: " + zonesResp.Status()})

		return
	}

	records := 0

	for _, z := range zones {
		var zd PdnsZoneDetails

		zr, zerr := httpc.R().
			SetContext(ctx).
			SetResult(&zd).
			Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, z.ID))
		if zerr != nil || zr.IsStatusFailure() {
			continue
		}

		for _, rr := range zd.RRsets {
			records += len(rr.Records)
		}
	}

	uptimeSec := getInt("uptime")
	resp := StatisticsResponse{
		Zones:      len(zones),
		Records:    records,
		Uptime:     humanUptime(uptimeSec),
		Status:     "online",
		UDPQueries: getInt("udp-queries"),
		TCPQueries: getInt("tcp-queries"),
		ServerID:   env.C.DnsServerId,
		StartedAt:  startedAtFromNow(uptimeSec).Format(time.RFC3339),
	}

	c.JSON(200, resp)
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
