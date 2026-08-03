package statistics

import "encoding/json"

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

type PdnsZone struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Serial     int64  `json:"serial"`
	URL        string `json:"url"`
	SoaEditApi string `json:"soa_edit_api"` //nolint:tagliatelle
}

type Stat struct {
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type ZoneDetails struct {
	Name   string  `json:"name"`
	RRsets []RRSet `json:"rrsets"`
}

type RRSet struct {
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	TTL     int        `json:"ttl"`
	Records []RRRecord `json:"records"`
}

type RRRecord struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
