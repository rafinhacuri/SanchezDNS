package zonas

type Soa struct {
	StartOfAuthority string `binding:"required"      json:"startOfAuthority"`
	Email            string `binding:"required"      json:"email"`
	Refresh          int    `binding:"required,gt=0" json:"refresh"`
	Retry            int    `binding:"required,gt=0" json:"retry"`
	Expire           int    `binding:"required,gt=0" json:"expire"`
	NegativeCacheTtl int    `binding:"required,gt=0" json:"negativeCacheTtl"`
}

type CreateZoneRequest struct {
	Domain string `binding:"required" json:"domain"`
	Soa    Soa    `binding:"required" json:"soa"`
	Type   string `binding:"required" json:"type"`
}

type ZoneFetch struct {
	Name   string `json:"name"`
	Serial int64  `json:"serial"`
	Nivel  string `json:"nivel"`
	Dnssec bool   `json:"dnssec"`
}

type ZonesResponse struct {
	Zones []ZoneFetch `json:"zones"`
}

type ZonePdns struct {
	Name   string `json:"name"`
	Serial int64  `json:"serial"`
	Dnssec *bool  `json:"dnssec,omitempty"`
}

type zoneDetail struct {
	Dnssec bool `json:"dnssec"`
}

//nolint:tagliatelle
type createZonePayload struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	SOAEditAPI string `json:"soa_edit_api"`
}

type cryptoKeyPayload struct {
	Active  bool   `json:"active"`
	KeyType string `json:"keytype"`
}
