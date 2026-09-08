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
	Dnssec     bool   `json:"dnssec"`
	Nsec3Param string `json:"nsec3param"`
}

type ZoneDetalhe struct {
	Dnssec bool
	Nsec3  bool
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

type nsec3ParamPayload struct {
	Nsec3Param string `json:"nsec3param"`
}

type CryptoKey struct {
	ID        int      `json:"id"`
	KeyType   string   `json:"keytype"`
	Active    bool     `json:"active"`
	Published bool     `json:"published"`
	Algorithm string   `json:"algorithm"`
	Bits      int      `json:"bits"`
	Ds        []string `json:"ds"`
}

type DnssecKey struct {
	KeyType   string   `json:"keytype"`
	Algorithm string   `json:"algorithm"`
	Bits      int      `json:"bits"`
	Active    bool     `json:"active"`
	Published bool     `json:"published"`
	Ds        []string `json:"ds"`
}

type DnssecStatus struct {
	Zone     string      `json:"zone"`
	Dnssec   bool        `json:"dnssec"`
	Nsec3    bool        `json:"nsec3"`
	DsStatus string      `json:"dsStatus"`
	DsPai    []string    `json:"dsPai"`
	Validado bool        `json:"validado"`
	Keys     []DnssecKey `json:"keys"`
}

type dohAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Data string `json:"data"`
}

//nolint:tagliatelle
type dohResponse struct {
	Status int         `json:"Status"`
	AD     bool        `json:"AD"`
	Answer []dohAnswer `json:"Answer"`
}
