package models

type PdnsZone struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Serial     int64  `json:"serial"`
	URL        string `json:"url"`
	SoaEditApi string `json:"soa_edit_api"`
}

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

type PdnsStat struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}
